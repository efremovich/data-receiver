package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/webapi"
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

	var notFoundElements int

	for _, meta := range saleReport {
		seller, err := s.getSeller(ctx, client.GetMarketPlace())
		if err != nil {
			return wrapErr(fmt.Errorf("ошибка получения данных о продавце %s модуль sales:%w", desc.Seller, err))
		}

		meta.Seller = seller
		// Проверим есть ли товар в базе, в случае отсутствия запросим его в 1с
		s2card, err := s.getSeller2Card(ctx, meta.Card.ExternalID, seller.ID)
		if err != nil {
			return wrapErr(fmt.Errorf("ошибка получения данных о продавце %s модуль sales reports:%w", desc.Seller, err))
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
		meta.Barcode, err = s.setBarcode(ctx, *barcode)
		if err != nil {
			return err
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

	if len(saleReport) > 0 {
		p := entity.PackageDescription{
			PackageType: entity.PackageTypeSaleReports,
			Seller:      desc.Seller,
			Delay:       desc.Delay,
			Cursor:      saleReport[len(saleReport)-1].ExternalID,
		}
		s.ReceiveSaleReport(ctx, p)
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
