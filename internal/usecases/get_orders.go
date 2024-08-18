package usecases

import (
	"context"
	"database/sql"
	"errors"

	"github.com/efremovich/data-receiver/internal/entity"
	aerror "github.com/efremovich/data-receiver/pkg/aerror"
)

func (s *receiverCoreServiceImpl) ReceiveOrders(ctx context.Context, desc entity.PackageDescription) aerror.AError {
	client := s.apiFetcher[desc.Seller]

	ordersMetaList, err := client.GetOrders(ctx, desc)
	if err != nil {
		return aerror.New(ctx, entity.GetDataFromExSources, err, "ошибка получение данные из внешнего источника %s в БД: %s ", "", err.Error())
	}

	attrs := make(map[string]interface{})
	attrs["количество данных"] = len(ordersMetaList)
	attrs["seller"] = desc.Seller

	var notFoundElements int
	for _, order := range ordersMetaList {
		wb2card, err := s.wb2cardrepo.SelectByNmid(ctx, order.Card.ExternalID)

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении wb2card %s в БД.", "wb")
		}
		// TODO В случае отсутствия в Wb2Card - добавлять в него
		if wb2card == nil {
			notFoundElements++
			continue
		}
		seller, err := s.sellerRepo.SelectByTitle(ctx, desc.Seller)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении seller %s в БД.", "wb")
		}

		warehouse, err := s.warehouserepo.SelectBySellerIDAndTitle(ctx, seller.ID, order.Warehouse.Title)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении warehouserepo %s в БД.", "wb")
		}

		orderData, err := s.orderrepo.SelectByCardIDAndDate(ctx, wb2card.CardID, desc.UpdatedAt)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении orderData %s в БД.", "wb")
		}
		if orderData == nil {
			_, err = s.orderrepo.Insert(ctx, entity.Order{
				ExternalID:   order.Order.ExternalID,
				Price:        order.Order.Price,
				Quantity:     1,
				Discount:     order.Order.Discount,
				SpecialPrice: order.Order.SpecialPrice,
				Status:       "",
				Type:         orderData.Type,
				Direction:    order.Order.Direction,
				CreatedAt:    desc.UpdatedAt,
				WarehouseID:  warehouse.ID,
				SellerID:     seller.ID,
				CardID:       wb2card.CardID,
			})
			if err != nil {
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении stock в БД.")
			}
		} else {
			orderData.Price = order.Order.Price
			orderData.Discount = order.Order.Discount
			orderData.SpecialPrice = order.Order.SpecialPrice
			err = s.orderrepo.UpdateExecOne(ctx, *orderData)
			if err != nil {
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении stock в БД.")
			}
		}

	}

	return nil
}
