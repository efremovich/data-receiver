package usecases

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/webapi"
	"github.com/efremovich/data-receiver/pkg/alogger"
)

func (s *receiverCoreServiceImpl) ReceiveSaleReport(ctx context.Context, desc entity.PackageDescription) error {
	clients := s.apiFetcher[desc.Seller]
	for _, client := range clients {
		err := s.receiveAndSaveSalesReport(ctx, client, desc)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *receiverCoreServiceImpl) receiveAndSaveSalesReport(ctx context.Context, client webapi.ExtAPIFetcher, desc entity.PackageDescription) error {
	saleReport, err := client.GetSaleReport(ctx, desc)
	if err != nil {
		return fmt.Errorf("ошибка получение данных о продажах из внешнего источника %s, %s", desc.Seller, err.Error())
	}

	alogger.InfoFromCtx(ctx, "Получение данных отчет по продажам за %s", desc.UpdatedAt.Format("02.01.2006"))

	var notFoundElements int

	seller, err := s.getSeller(ctx, client.GetMarketPlace())
	if err != nil {
		return wrapErr(fmt.Errorf("ошибка получения данных о продавце %s модуль sales:%w", desc.Seller, err))
	}

	for _, meta := range saleReport {
		meta.Seller = seller
		// Проверим есть ли товар в базе, в случае отсутствия запросим его в 1с
		s2card, err := s.getSeller2Card(ctx, meta.Card.ExternalID, seller.ID)
		if err != nil {
			alogger.InfoFromCtx(ctx, "ошибка получения данных о продавце %s модуль sales reports: %s", desc.Seller, err.Error())
			// return wrapErr(fmt.Errorf("ошибка получения данных о продавце %s модуль sales reports:%w", desc.Seller, err))
			notFoundElements++
			continue
		}

		card, err := s.getCardByID(ctx, s2card.CardID)
		if errors.Is(err, ErrObjectNotFound) {
			notFoundElements++
			continue
		} else if err != nil {
			return err
		}

		meta.Card = card

		// Size
		size, err := s.getSizeByTitle(ctx, meta.Size.TechSize)
		if err != nil {
			return err
		}

		meta.Size = size

		meta.Warehouse.SellerID = seller.ID
		warehouse, err := s.setWarehouse(ctx, meta.Warehouse)
		if err != nil {
			return err
		}

		meta.Warehouse = warehouse

		if meta.Pvz != nil {
			meta.Pvz, err = s.setPvz(ctx, meta.Pvz)
			if err != nil {
				return err
			}
		}

		barcode := meta.Barcode
		_, err = s.setBarcode(ctx, *barcode)
		if err != nil {
			fmt.Println(err.Error())
			// return err
		}

		meta.Order, err = s.getOrderByExternalID(ctx, meta.Order.ExternalID)
		if err != nil {
			return err
		}

		err = s.setSaleReport(ctx, &meta)
		if err != nil {
			return err
		}
	}

	alogger.InfoFromCtx(ctx, "Загружена информация по отчету продаж всего: %d из них не найдено %d", len(saleReport), notFoundElements)

	alogger.InfoFromCtx(ctx, "постановка задачи в очередь %d", desc.Limit)

	if desc.Limit > 0 {

		p := entity.PackageDescription{
			PackageType: entity.PackageTypeSaleReports,
			Seller:      desc.Seller,
			Delay:       desc.Delay,
			UpdatedAt:   desc.UpdatedAt.Add(-24 * time.Hour),
			// Cursor:      saleReport[len(saleReport)-1].ExternalID,
			Limit: desc.Limit - 1,
		}
		err := s.ReceiveSaleReport(ctx, p)
		if err != nil {
			return fmt.Errorf("ошибка постановки задачи в очередь: %w", err)
		}

		alogger.InfoFromCtx(ctx, "Создана очередь для получения отчета по продажам на %s", p.UpdatedAt.Format("02.01.2006"))
	} else {
		alogger.InfoFromCtx(ctx, "Все элементы обработаны")
	}
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
