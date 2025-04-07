package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/webapi"
	"github.com/efremovich/data-receiver/pkg/alogger"
	"github.com/efremovich/data-receiver/pkg/logger"
	"golang.org/x/sync/errgroup"
)

func (s *receiverCoreServiceImpl) ReceiveStocks(ctx context.Context, desc entity.PackageDescription) error {
	clients := s.apiFetcher[desc.Seller]

	g, gCtx := errgroup.WithContext(ctx)

	for _, c := range clients {
		client := c

		g.Go(func() error {
			return s.receiveAndSaveStocks(gCtx, client, desc)
		})
	}

	// Ждем завершения всех горутин и проверяем наличие ошибок
	if err := g.Wait(); err != nil {
		if errors.Is(err, context.Canceled) {
			alogger.WarnFromCtx(ctx, "Операция была отменена: %v", err)
			return nil
		}
		return fmt.Errorf("ошибка при обработке клиентов: %w", err)
	}

	alogger.InfoFromCtx(ctx, "постановка задачи в очередь %d", desc.Limit)

	if desc.Limit > 0 {
		p := entity.PackageDescription{
			PackageType: entity.PackageTypeStock,
			UpdatedAt:   desc.UpdatedAt.Add(-24 * time.Hour),
			Limit:       desc.Limit - 1,
			Seller:      desc.Seller,
		}

		err := s.brokerPublisher.SendPackage(ctx, &p)
		if err != nil {
			return fmt.Errorf("ошибка постановки задачи в очередь: %w", err)
		}

		alogger.InfoFromCtx(ctx, "Создана очередь на получение остатков на дату: %s", p.UpdatedAt.Format("02.01.2006"))
	} else {
		alogger.InfoFromCtx(ctx, "Все элементы обработаны")
	}
	return nil
}

func (s *receiverCoreServiceImpl) receiveAndSaveStocks(ctx context.Context, client webapi.ExtAPIFetcher, desc entity.PackageDescription) error {
	var notFoundElements int

	stockMetaList, err := client.GetStocks(ctx, desc)
	if err != nil {
		return fmt.Errorf("ошибка получение данные из внешнего источника %s, %w", desc.Seller, err)
	}

	alogger.InfoFromCtx(ctx, "Получение данных об остатках за %s", desc.UpdatedAt.Format("02.01.2006"))

	seller, err := s.getSeller(ctx, client.GetMarketPlace())
	if err != nil {
		return wrapErr(fmt.Errorf("ошибка получения данных о продавце %s модуль stocks: %w", desc.Seller, err))
	}

	for _, meta := range stockMetaList {
		// Проверим есть ли товар в базе, в случае отсутствия запросим его в 1с
		_, err = s.getSeller2Card(ctx, meta.Seller2Card.ExternalID, seller.ID)
		// Получаем ошибку что такой записи нет, поищем карточку товара в 1с
		if errors.Is(err, ErrObjectNotFound) {
			query := make(map[string]string)
			query["barcode"] = meta.Barcode.Barcode
			query["article"] = meta.Card.VendorID

			descForOdin := entity.PackageDescription{
				Seller: "odinc",
				Query:  query,
			}

			err := s.ReceiveCards(ctx, descForOdin)
			if err != nil {
				return err
			}
		} else if err != nil {
			logger.GetLoggerFromContext(ctx).Errorf("ошибка получения данных о товаре %s модуль stocks: %s", desc.Seller, err.Error())
			continue
		}

		card, err := s.getCardByVendorID(ctx, meta.Card.VendorID)
		if errors.Is(err, ErrObjectNotFound) {
			// Нам не удалось получить запись, значит данные по этому товару исчезли, пропускаем загрузку остатков
			// TODO Писать эти данные в jaeger
			notFoundElements++
			continue
		} else if err != nil {
			return err
		}
		// Проверим и создадим связь продавца и товара
		meta.Seller2Card.CardID = card.ID
		meta.Seller2Card.SellerID = seller.ID

		_, err = s.setSeller2Card(ctx, meta.Seller2Card)
		if err != nil {
			return err
		}

		// Size
		size, err := s.setSize(ctx, &meta.Size)
		if err != nil {
			return err
		}

		// Добавление данных
		// PriceSize
		meta.PriceSize.CardID = card.ID
		meta.PriceSize.SizeID = size.ID

		priceSize, err := s.setPriceSize(ctx, meta.PriceSize)
		if err != nil {
			return err
		}

		// Barcode
		meta.Barcode.SellerID = seller.ID
		meta.Barcode.PriceSizeID = priceSize.ID

		barcode, err := s.setBarcode(ctx, meta.Barcode)
		if err != nil {
			return err
		}

		// Warehouse
		meta.Warehouse.SellerID = seller.ID

		warehouse, err := s.setWarehouse(ctx, &meta.Warehouse)
		if err != nil {
			return wrapErr(fmt.Errorf("ошибка получения данных по складам хранения модуль stock:%w", err))
		}

		stock := meta.Stock
		stock.BarcodeID = barcode.ID
		stock.WarehouseID = warehouse.ID
		stock.SellerID = seller.ID
		stock.CardID = card.ID
		stock.UpdatedAt = time.Now()

		_, err = s.setStock(ctx, stock)
		if err != nil {
			return err
		}
	}

	alogger.InfoFromCtx(ctx, "Загружена информация о остатке всего: %d из них не найдено %d", len(stockMetaList), notFoundElements)

	return nil
}

func (s *receiverCoreServiceImpl) setStock(ctx context.Context, in entity.Stock) (*entity.Stock, error) {
	stock, err := s.stockrepo.SelectByBarcode(ctx, in.BarcodeID, in.CreatedAt)

	if errors.Is(err, ErrObjectNotFound) {
		stock, err = s.stockrepo.Insert(ctx, in)
	}

	if err != nil {
		return nil, wrapErr(fmt.Errorf("ошибка получения и обновления остатка: %w", err))
	}
	return stock, nil
}
