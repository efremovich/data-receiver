package usecases

import (
	"context"
	"errors"
	"fmt"

	"github.com/efremovich/data-receiver/internal/entity"
)

func (s *receiverCoreServiceImpl) ReceiveWarehouses(ctx context.Context, desc entity.PackageDescription) error {
	client := s.apiFetcher[desc.Seller]
	warehouses, err := client.GetWarehouses(ctx)
	if err != nil {
		return wrapErr(fmt.Errorf("ошбка при получении данных из источника %s : %w", desc.Seller, err))
	}

	for _, in := range warehouses {
		_, err := s.setWarehouse(ctx, &in)
		if err != nil {
			return err
		}

	}
	return nil
}

func (s *receiverCoreServiceImpl) setWarehouse(ctx context.Context, in *entity.Warehouse) (*entity.Warehouse, error) {
	warehouse, err := s.warehouserepo.SelectBySellerIDAndTitle(ctx, in.SellerID, in.Title)
	if errors.Is(err, ErrObjectNotFound) {
		warehouseType, err := s.getWarehouseType(ctx, in.TypeName)
		if err != nil {
			return nil, err
		}

		warehouse, err = s.warehouserepo.Insert(ctx, entity.Warehouse{
			ExternalID: in.ExternalID,
			Title:      in.Title,
			Address:    in.Address,
			TypeID:     warehouseType.ID,
			SellerID:   in.SellerID,
		})
		if err != nil {
			return nil, wrapErr(fmt.Errorf("Ошибка при получении данных: %w", err))
		}

	}

	return warehouse, nil
}

func (s *receiverCoreServiceImpl) getWarehouseType(ctx context.Context, typeTitle string) (*entity.WarehouseType, error) {
	warehouseType, err := s.warehousetyperepo.SelectByTitle(ctx, typeTitle)

	if errors.Is(err, ErrObjectNotFound) {
		warehouseType, err = s.warehousetyperepo.Insert(ctx, entity.WarehouseType{
			Title: typeTitle,
		})
	}
	if err != nil {
		return nil, wrapErr(fmt.Errorf("Ошибка при получении данных: %w", err))
	}
	return warehouseType, nil
}