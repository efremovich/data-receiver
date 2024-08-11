package usecases

import (
	"context"
	"database/sql"
	"errors"

	"github.com/efremovich/data-receiver/internal/entity"
	aerror "github.com/efremovich/data-receiver/pkg/aerror"
)

func (s *receiverCoreServiceImpl) ReceiveWarehouses(ctx context.Context) aerror.AError {
	client := s.apiFetcher["wb"]
	warehouses, err := client.GetWarehouses(ctx)
	if err != nil {
		return aerror.New(ctx, entity.GetDataFromExSources, err, "Ошибка при получении warehouse %s", "wb")
	}

	seller, err := s.sellerRepo.SelectByTitle(ctx, "wb")
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении Seller %s в БД.", "wb")
	}

	if seller == nil {
		seller, err = s.sellerRepo.Insert(ctx, entity.Seller{
			Title:     "wb",
			IsEnabled: true,
		})
		if err != nil {
			return aerror.New(ctx, entity.InsertDataErrorID, err, "Ошибка при сохранении Seller товара %d в БД.", seller.ID)
		}
	}

	if err != nil {
		return aerror.New(ctx, entity.GetDataFromExSources, err, "ошибка получение данные из внешнего источника %s в БД: %s ", "", err.Error())
	}

	for _, elem := range warehouses {
		warehouse, err := s.warehousetyperepo.SelectByTitle(ctx, elem.Title)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении Warehouese %s в БД.", "wb")
		}

		if warehouse == nil {
			_, err := s.warehouserepo.Insert(ctx, entity.Warehouse{
				ExternalID: elem.ExternalID,
				Title:      elem.Title,
				Address:    elem.Address,
				TypeID:     elem.TypeID,
				SellerID:   seller.ID,
			})
			if err != nil {
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении warehouse %s в БД.", elem.Title)
			}
		}

	}
	return nil
}
