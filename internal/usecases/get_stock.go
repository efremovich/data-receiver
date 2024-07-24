package usecases

import (
	"context"
	"database/sql"
	"errors"

	"github.com/efremovich/data-receiver/internal/entity"
	aerror "github.com/efremovich/data-receiver/pkg/aerror"
)

func (s *receiverCoreServiceImpl) ReceiveStocks(ctx context.Context, desc entity.PackageDescription) aerror.AError {
	client := s.apiFetcher["wb"]
	stockMetaList, err := client.GetStocks(ctx, desc)
	if err != nil {
		return aerror.New(ctx, entity.GetDataFromExSources, err, "ошибка получение данные из внешнего источника %s в БД: %s ", "", err.Error())
	}

	for _, stock := range stockMetaList {
		wb2card, err := s.wb2cardrepo.SelectByNmid(ctx, stock.Wb2Card.NMID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении wb2card %s в БД.", "wb")
		}

		card, err := s.cardCharRepo.SelectByID(ctx, wb2card.CardID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении card %s в БД.", "wb")
		}

		seller, err := s.sellerRepo.SelectByTitle(ctx, "wb")

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении card %s в БД.", "wb")
		}
	}
}
