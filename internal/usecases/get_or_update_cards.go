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
	client := s.apiFetcher[desc.Seller]

	cards, err := client.GetCards(ctx, desc)
	if err != nil {
		return aerror.New(ctx, entity.GetDataFromExSources, err, "ошибка получение данные из внешнего источника %s в БД: %s ", "", err.Error())
	}

	for _, card := range cards {
		// Seller
		seller, err := s.getSeller(ctx, desc)
		if err != nil {
			return err
		}

		// Brands
		brand, err := s.getBrand(ctx, card.Brand, seller)
		if err != nil {
			return err
		}

		card.Brand = *brand
		// Wb2Card
		err = s.setWb2Card(ctx, &card)
		if err != nil {
			return err
		}

		// Characteristics
		err = s.setCharacterisitc(ctx, card)
		if err != nil {
			return err
		}

		// Sizes
		err = s.setSizes(ctx, card)
		if err != nil {
			return err
		}

		// Dimensions
		err = s.setDimension(ctx, card)
		if err != nil {
			return err
		}

		// Mediafile
		err = s.setMediaFile(ctx, card)
		if err != nil {
			return err
		}

		// Categorites
		err = s.setCategory(ctx, card, seller)
		if err != nil {
			return err
		}
	}

	if len(cards) == desc.Limit {
		p := entity.PackageDescription{
			PackageType: entity.PackageTypeCard,
			Cursor:      int(cards[len(cards)-1].ExternalID),
			UpdatedAt:   cards[len(cards)-1].UpdatedAt,
			Limit:       desc.Limit,
			Seller:      desc.Seller,
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

func (s *receiverCoreServiceImpl) setCategory(ctx context.Context, card entity.Card, seller *entity.Seller) aerror.AError {
	for _, cat := range card.Categories {
		category, err := s.categoryRepo.SelectByTitle(ctx, cat.Title)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении category %s в БД.", cat.Title)
		}

		if category == nil {
			_, err = s.categoryRepo.Insert(ctx, entity.Category{
				Title:      cat.Title,
				ExternalID: cat.ExternalID,
				CardID:     card.ID,
				SellerID:   seller.ID,
				ParentID:   0,
			})
			if err != nil {
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении category %s в БД.", cat.Title)
			}
		}
	}

	return nil
}

func (s *receiverCoreServiceImpl) setMediaFile(ctx context.Context, card entity.Card) aerror.AError {
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
				return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении dimension %d в БД.", card.ID)
			}
		}
	}

	return nil
}

func (s *receiverCoreServiceImpl) setDimension(ctx context.Context, card entity.Card) aerror.AError {
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

	return nil
}

func (s *receiverCoreServiceImpl) setSizes(ctx context.Context, card entity.Card) aerror.AError {
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

	return nil
}

func (s *receiverCoreServiceImpl) setCharacterisitc(ctx context.Context, card entity.Card) aerror.AError {
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

	return nil
}

func (s *receiverCoreServiceImpl) setWb2Card(ctx context.Context, card *entity.Card) aerror.AError {
	wb2card, err := s.wb2cardrepo.SelectByNmid(ctx, card.ExternalID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении Wb2Card %s в БД.", card.Title)
	}

	if wb2card == nil {
		newCard, err := s.cardRepo.Insert(ctx, *card)
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
		err = s.cardRepo.UpdateExecOne(ctx, *card)

		if err != nil {
			return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении карточки товара %s в БД.", card.Title)
		}
	}

	return nil
}

func (s *receiverCoreServiceImpl) getBrand(ctx context.Context, brandIn entity.Brand, seller *entity.Seller) (*entity.Brand, aerror.AError) {
	brand, err := s.brandRepo.SelectByTitleAndSeller(ctx, brandIn.Title, seller.ID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, aerror.New(ctx, entity.SelectPkgErrorID, err, "Ошибка при получении Brand %s в БД.", brandIn.Title)
	}

	if brand == nil {
		brand, err = s.brandRepo.Insert(ctx, entity.Brand{
			ExternalID: brandIn.ExternalID,
			Title:      brandIn.Title,
			SellerID:   seller.ID,
		})
		if err != nil {
			return nil, aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении карточки товара %s в БД.", brandIn.Title)
		}
	}

	return brand, nil
}

func (s *receiverCoreServiceImpl) getSeller(ctx context.Context, desc entity.PackageDescription) (*entity.Seller, aerror.AError) {
	seller, err := s.sellerRepo.SelectByTitle(ctx, desc.Seller)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, aerror.New(ctx, entity.SelectDataErrorID, err, "Ошибка при получении Seller %s в БД.", "wb")
	}

	if seller == nil {
		seller, err = s.sellerRepo.Insert(ctx, entity.Seller{
			Title:     desc.Seller,
			IsEnabled: true,
		})
		if err != nil {
			return nil, aerror.New(ctx, entity.InsertDataErrorID, err, "Ошибка при сохранении Seller: %s в БД.", desc.Seller)
		}
	}

	return seller, nil
}
