package usecases

import (
	"context"
	"database/sql"
	"errors"

	"github.com/efremovich/data-receiver/internal/entity"
	aerror "github.com/efremovich/data-receiver/pkg/aerror"
)

func (s *receiverCoreServiceImpl) saveBrandToDB(ctx context.Context, in *entity.Brand, sellerID int64) error {
	// Brands
	brand, err := s.brandRepo.SelectByTitleAndSeller(ctx, in.Title, sellerID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return aerror.New(ctx, entity.SelectPkgErrorID, err, "Ошибка при получении Brand %s в БД.", brand.Title)
	}

	if brand == nil {
		brand, err = s.brandRepo.Insert(ctx, entity.Brand{
			ExternalID: in.ExternalID,
			Title:      in.Title,
			SellerID:   sellerID,
		})
		if err != nil {
			return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении карточки товара %s в БД.", brand.Title)
		}
	}

	return nil
}
