package usecases

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/efremovich/data-receiver/internal/entity"
	aerror "github.com/efremovich/data-receiver/pkg/aerror"
)

func (s *receiverCoreServiceImpl) ReceiveCards(ctx context.Context, desc entity.PackageDescription) aerror.AError {
	client := s.apiFetcher["wb"]
	cards, err := client.GetCards(ctx, desc)
	if err != nil {
		return aerror.New(ctx, entity.GetDataFromExSources, err, "ошибка получение данные из внешнего источника %s в БД: %s ", "", err.Error())
	}
	for _, card := range cards {
		// Seller
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
				return aerror.New(ctx, entity.InsertDataErrorID, err, "Ошибка при сохранении Seller товара %s в БД.", card.Title)
			}
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
		wb2card, err := s.wb2cardrepo.SelectByNmid(ctx, card.ExternalID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении Wb2Card %s в БД.", card.Title)
		}
		if wb2card == nil {
			newCard, err := s.cardRepo.Insert(ctx, card)
			if err != nil {
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении карточки товара %s в БД.", card.Title)
			}
			card.ID = newCard.ID
			_, err = s.wb2cardrepo.Insert(ctx, entity.Wb2Card{
				NMID:   card.ExternalID,
				KTID:   0,
				NMUUID: "",
				CardID: card.ID,
			})
			if err != nil {
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении карточки товара %s в БД.", card.Title)
			}
		} else {
			card.ID = wb2card.CardID
			err = s.cardRepo.UpdateExecOne(ctx, card)
			if err != nil {
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении карточки товара %s в БД.", card.Title)
			}
		}

		// Characteristics
		for _, elem := range card.Characteristics {
			char, err := s.charRepo.SelectByTitle(ctx, elem.Title)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return aerror.New(ctx, entity.SelectPkgErrorID, err, "Ошибка при получении Charactristic %s в БД.", card.Title)
			}
			if char == nil {
				char, err = s.charRepo.Insert(ctx, entity.Characteristic{
					Title: elem.Title,
				})
				if err != nil {
					return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении Charactristic %s в БД.", card.Title)
				}
			}
			cardChar, err := s.cardCharRepo.SelectByCardIDAndCharID(ctx, card.ID, char.ID)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении CardCharactristic %s в БД.", card.Title)
			}

			if cardChar == nil {
				_, err = s.cardCharRepo.Insert(ctx, entity.CardCharacteristic{
					Value:            elem.Value,
					Title:            elem.Title,
					CharacteristicID: char.ID,
					CardID:           card.ID,
				})
				if err != nil {
					return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении CardCharactristic %s в БД.", card.Title)
				}
			}
		}
		// Sizes
		for _, elem := range card.Sizes {
			size, err := s.sizerepo.SelectByTitle(ctx, elem.Title)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении Size %s в БД.", card.Title)
			}
			if size == nil {
				_, err = s.sizerepo.Insert(ctx, entity.Size{
					ExternalID: elem.ExternalID,
					TechSize:   elem.TechSize,
					Title:      elem.Title,
				})
				if err != nil {
					return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении Size %s в БД.", card.Title)
				}
			}
		}
		// Dimensions
		dimension, err := s.dimensionrepo.SelectByCardID(ctx, card.ID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении dimension %d в БД.", card.ID)
		}
		if dimension == nil {
			_, err = s.dimensionrepo.Insert(ctx, entity.Dimension{
				Width:   card.Dimension.Width,
				Height:  card.Dimension.Height,
				Length:  card.Dimension.Length,
				IsVaild: card.Dimension.IsVaild,
				CardID:  card.ID,
			})
			if err != nil {
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении dimension %d в БД.", card.ID)
			}
		}
		// Mediafile
		for _, elem := range card.MediaFile {
			mf, err := s.mediafilerepo.SelectByCardID(ctx, card.ID, elem.Link)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении dimension %d в БД.", card.ID)
			}
			if mf == nil {
				_, err = s.mediafilerepo.Insert(ctx, entity.MediaFile{
					Link:   elem.Link,
					TypeID: elem.TypeID,
					CardID: card.ID,
				})
				if err != nil {
					return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении MediaFile %d в БД.", card.ID)
				}
			}
		}
	}

	if len(cards) == desc.Limit {
		p := entity.PackageDescription{
			PackageType: "CARD",
			Cursor:      int(cards[len(cards)-1].ExternalID),
			UpdatedAt:   &cards[len(cards)-1].UpdatedAt,
			Limit:       desc.Limit,
		}
		err = s.brokerPublisher.SendPackage(ctx, &p)
		if err != nil {
			return aerror.New(ctx, entity.BrokerSendErrorID, err, "Ошибка постановки задачи в очередь")
		}
	} else {
		fmt.Println("its all")
	}
	return nil
}
