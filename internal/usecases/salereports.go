package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/webapi"
	"github.com/efremovich/data-receiver/pkg/alogger"
)

func (s *receiverCoreServiceImpl) ReceiveSaleReport(ctx context.Context, desc entity.PackageDescription) error {
	clients := s.apiFetcher[desc.Seller]

	g, gCtx := errgroup.WithContext(ctx)

	for _, c := range clients {
		client := c

		g.Go(func() error {
			return s.receiveAndSaveSalesReport(gCtx, client, desc)
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
			PackageType: entity.PackageTypeSaleReports,
			Seller:      desc.Seller,
			Delay:       desc.Delay,
			UpdatedAt:   desc.UpdatedAt.Add(-24 * time.Hour),
			Limit:       desc.Limit - 1,
		}
		err := s.brokerPublisher.SendPackage(ctx, &p)
		if err != nil {
			return fmt.Errorf("ошибка постановки задачи в очередь: %w", err)
		}

		alogger.InfoFromCtx(ctx, "Создана очередь для получения отчета по продажам на %s", p.UpdatedAt.Format("02.01.2006"))
	} else {
		alogger.InfoFromCtx(ctx, "Все элементы обработаны")
	}

	return nil
}

func (s *receiverCoreServiceImpl) receiveAndSaveSalesReport(ctx context.Context, client webapi.ExtAPIFetcher, desc entity.PackageDescription) error {
	saleReport, err := client.GetSaleReport(ctx, desc)
	if err != nil {
		return fmt.Errorf("ошибка получение данных о продажах из внешнего источника %s, %s", desc.Seller, err.Error())
	}

	var notFoundElements int

	seller, err := s.getSeller(ctx, client.GetMarketPlace())
	if err != nil {
		return wrapErr(fmt.Errorf("ошибка получения данных о продавце %s модуль sales:%w", desc.Seller, err))
	}

	alogger.InfoFromCtx(ctx, "Получение данных отчет по продажам за %s маркетплейс %s", desc.UpdatedAt.Format("02.01.2006"), seller.Title)

	for _, meta := range saleReport {
		meta.Seller = seller
		// Проверим есть ли товар в базе, в случае отсутствия запросим его в 1с
		s2card, err := s.getSeller2Card(ctx, meta.Card.ExternalID, seller.ID)
		if err != nil {
			alogger.InfoFromCtx(ctx, "ошибка получения данных о продавце %s модуль sales reports: %s", desc.Seller, err.Error())
			notFoundElements++

			continue
		}

		card, err := s.getCardByID(ctx, s2card.CardID)
		if errors.Is(err, ErrObjectNotFound) {
			notFoundElements++
			continue
		} else if err != nil {
			alogger.InfoFromCtx(ctx, "ошибка получения данных о карточки товара %s модуль sales reports: %s", desc.Seller, err.Error())
			return err
		}

		meta.Card = card

		// Size
		size, err := s.setSize(ctx, meta.Size)
		if err != nil {
			alogger.InfoFromCtx(ctx, "ошибка получения данных о размер %s модуль sales reports: %s", desc.Seller, err.Error())
			return err
		}

		meta.Size = size

		meta.Warehouse.SellerID = seller.ID
		warehouse, err := s.setWarehouse(ctx, meta.Warehouse)
		if err != nil {
			alogger.InfoFromCtx(ctx, "ошибка получения данных о складах %s модуль sales reports: %s", desc.Seller, err.Error())
			return err
		}

		meta.Warehouse = warehouse

		if meta.Pvz != nil {
			meta.Pvz, err = s.setPvz(ctx, meta.Pvz)
			if err != nil {
				alogger.InfoFromCtx(ctx, "ошибка получения данных о ПВЗ %s модуль sales reports: %s", desc.Seller, err.Error())
				return err
			}
		}

		meta.Order, err = s.getOrderByExternalID(ctx, meta.Order.ExternalID)
		if err != nil {
			alogger.InfoFromCtx(ctx, "ошибка получения данных о заказе %s модуль sales reports: %s", desc.Seller, err.Error())
			return err
		}

		err = s.setSaleReport(ctx, &meta)
		if err != nil {
			alogger.InfoFromCtx(ctx, "ошибка получения данных о отчетах по продажам %s модуль sales reports: %s", desc.Seller, err.Error())
			return err
		}
	}

	alogger.InfoFromCtx(ctx, "Загружена информация по отчету продаж всего: %d из них не найдено %d", len(saleReport), notFoundElements)

	return nil
}

func (s *receiverCoreServiceImpl) setSaleReport(ctx context.Context, in *entity.SaleReport) error {
	_, err := s.saleReportRepo.SelectByExternalID(ctx, in.ExternalID, in.SaleDate)
	if errors.Is(err, ErrObjectNotFound) {
		_, err = s.saleReportRepo.Insert(ctx, *in)
	}

	if err != nil {
		return err
	}
	return nil
}
