package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/pkg/alogger"
)

func (s *receiverCoreServiceImpl) ReceiveStocks(ctx context.Context, desc entity.PackageDescription) error {
	client := s.apiFetcher[desc.Seller]

	stockMetaList, err := client.GetStocks(ctx, desc)
	if err != nil {
		return fmt.Errorf("ошибка получение данные из внешнего источника %s, %s", desc.Seller, err)
	}

	alogger.InfoFromCtx(ctx, "Получение данных об остатках за %s", desc.UpdatedAt.Format("02.01.2006"))

	var notFoundElements int

	for _, meta := range stockMetaList {
		seller, err := s.getSeller(ctx, desc.Seller)
		if err != nil {
			return wrapErr(fmt.Errorf("ошибка получения данных о продавце %s модуль stocks: %w", desc.Seller, err))
		}

		// Проверим есть ли товар в базе, в случае отсутствия запросим его в 1с
		_, err = s.getSeller2Card(ctx, meta.Seller2Card.ExternalID, seller.ID)
		// Получаем ошибку что такой записи нет, поищем карточку товара в 1с
		if errors.Is(err, ErrObjectNotFound) {
			query := make(map[string]string)
			query["barcode"] = meta.Barcode.Barcode
			query["article"] = meta.SupplierArticle

			desc := entity.PackageDescription{
				Seller: "1c",
				Query:  query,
			}
			err := s.ReceiveCards(ctx, desc)
			if err != nil {
				return err
			}
		}
		if err != nil {
			return wrapErr(fmt.Errorf("ошибка получения данных отсутствует связь между продавцом и товаром модуль stocks:%w", desc.Seller, err))
		}

		card, err := s.getCardByVendorCode(ctx, meta.SupplierArticle)
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
		size, err := s.getSizeByTitle(ctx, meta.Size.TechSize)
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

	if desc.Limit > 0 {
		p := entity.PackageDescription{
			PackageType: entity.PackageTypeStock,

			UpdatedAt: desc.UpdatedAt.Add(-24 * time.Hour),
			Limit:     desc.Limit - 1,
			Seller:    desc.Seller,
		}

		err = s.brokerPublisher.SendPackage(ctx, &p)
		if err != nil {
			return fmt.Errorf("Ошибка постановки задачи в очередь: %w", err)
		}

		alogger.InfoFromCtx(ctx, "Создана очередь на получение остатков на дату: %s", p.UpdatedAt.Format("02.01.2006"))
	} else {
		alogger.InfoFromCtx(ctx, "Все элементы обработаны")
	}

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