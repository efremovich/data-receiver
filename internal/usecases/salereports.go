package usecases

import (
	"context"
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
	_, err := client.GetSaleReport(ctx, desc)
	if err != nil {
		return fmt.Errorf("ошибка получение данных о продажах из внешнего источника %s, %w", desc.Seller, err)
	}
	return nil
}
