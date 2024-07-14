package usecases

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/efremovich/data-receiver/internal/entity"
	aerror "github.com/efremovich/data-receiver/pkg/aerror"
)

func (s *receiverCoreServiceImpl) ReceiveCards(ctx context.Context, cursor int) aerror.AError {
	client := s.apiFetcher["wb"]
	cards, cursor, err := client.GetCards(ctx, cursor)
	if err != nil {
		return aerror.New(ctx, entity.GetDataFromExSources, err, "ошибка получение данные из внешнего источника %s в БД: %s ", "", err.Error())
	}

	for _, card := range cards {
		// Seller
		seller, err := s.sellerRepo.SelectByTitle(ctx, "wb")
		if err != nil {
			return aerror.New(ctx, entity.SelectPkgErrorID, err, "Ошибка при получении Seller %s в БД.", "wb")
		}

		// Brands
		brand, err := s.brandRepo.SelectByTitleAndSeller(ctx, card.Brand.Title, seller.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectPkgErrorID, err, "Ошибка при получении Brand %s в БД.", card.Title)
		}
		if brand == nil {

			brand, err = s.brandRepo.Insert(ctx, entity.Brand{
				ExternalID: card.Brand.ExternalID,
				Title:      card.Brand.Title,
				SellerID:   seller.ID,
			})
			if err != nil {
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении карточки товара %s в БД.", card.Title)
			}

		}

		card.Brand = *brand

		// Card
		upCard, err := s.cardRepo.Insert(ctx, card)
		if err != nil {
			return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении карточки товара %s в БД.", card.Title)
		}

		// Characteristics
		for _, elem := range card.Characteristics {
			char, err := s.charRepo.SelectByTitle(ctx, elem.Title)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return aerror.New(ctx, entity.SelectPkgErrorID, err, "Ошибка при получении CardCharactristic %s в БД.", card.Title)
			}
			if char == nil {
				char, err = s.charRepo.Insert(ctx, entity.Characteristic{
					Title: elem.Title,
				})
				if err != nil {
					return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении CardCharactristic %s в БД.", card.Title)
				}
			}

			cardChar, err := s.charRepo.SelectByCardID(ctx, upCard.ID, strings.Join(elem.Value, ","))
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return aerror.New(ctx, entity.SelectPkgErrorID, err, "Ошибка при получении CardCharactristic %s в БД.", card.Title)
			}
			if cardChar == nil {
				addedCardChar, err := s.charRepo.InsertCardChar(ctx, entity.CardCharacteristic{
					Value:            elem.Value,
					Title:            elem.Title,
					CharacteristicID: char.ID,
					CardID:           upCard.ID,
				})
				if err != nil {
					return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении CardCharactristic %s в БД.", card.Title)
				}
				cardChar = append(cardChar, addedCardChar)
			}
		}
	}
	p := entity.PackageDescription{
		PackageType: "CARD",
		Cursor:      cursor,
	}
	err = s.brokerPublisher.SendPackage(ctx, &p)
	if err != nil {
		return aerror.New(ctx, entity.BrokerSendErrorID, err, "Ошибка постановки задачи в очередь")
	}
	return nil
}
