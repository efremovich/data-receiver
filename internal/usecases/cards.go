package usecases

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/webapi"
	"github.com/efremovich/data-receiver/pkg/alogger"
)

func (s *receiverCoreServiceImpl) ReceiveCards(ctx context.Context, desc entity.PackageDescription) error {
	clients := s.apiFetcher[desc.Seller]

	for _, client := range clients {
		err := s.receiveAndSaveCard(ctx, client, desc)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *receiverCoreServiceImpl) receiveAndSaveCard(ctx context.Context, client webapi.ExtAPIFetcher, desc entity.PackageDescription) error {
	cards, err := client.GetCards(ctx, desc)
	if err != nil {
		return fmt.Errorf("ошибка получение данные из внешнего источника %s, %w", desc.Seller, err)
	}

	alogger.InfoFromCtx(ctx, "Начали загрузку карточек товара %d от %s", len(cards), desc.Seller)

	for _, in := range cards {
		s.metricsCollector.IncServiceDocsTaskCounter()
		// Seller
		seller, err := s.getSeller(ctx, desc.Seller)
		if err != nil {
			return err
		}

		// Brands
		brand, err := s.getBrand(ctx, in.Brand, seller)
		if err != nil {
			return err
		}

		in.Brand = *brand

		// Cards
		card, err := s.setCard(ctx, in)
		if err != nil {
			return err
		}

		// Seller2Card
		seller2card := entity.Seller2Card{
			ExternalID: in.ExternalID,
			CardID:     card.ID,
			SellerID:   seller.ID,
		}

		_, err = s.setSeller2Card(ctx, seller2card)
		if err != nil {
			return err
		}

		// Characteristics
		_, err = s.setCharacterisitc(ctx, card)
		if err != nil {
			return err
		}

		// Sizes
		for _, size := range card.Sizes {
			_, err = s.setSize(ctx, size)
			if err != nil {
				return err
			}
		}

		// Dimensions
		_, err = s.setDimension(ctx, card)
		if err != nil {
			return err
		}

		// Mediafile
		card.MediaFile = in.MediaFile
		_, err = s.setMediaFile(ctx, card)
		if err != nil {
			return err
		}

		// Categorites
		categories, err := s.setCategory(ctx, card, seller)
		if err != nil {
			return err
		}

		err = s.setCardCategories(ctx, card.ID, categories)
		if err != nil {
			return err
		}
	}

	if len(cards) != desc.Limit {
		alogger.InfoFromCtx(ctx, "Задание успешно завершено")
	} else if len(cards) > 0 {
		lastID := strconv.Itoa(int(cards[len(cards)-1].ExternalID))
		p := entity.PackageDescription{
			PackageType: entity.PackageTypeCard,
			Cursor:      lastID,
			UpdatedAt:   cards[len(cards)-1].UpdatedAt,
			Limit:       desc.Limit,
			Seller:      desc.Seller,
		}

		err = s.brokerPublisher.SendPackage(ctx, &p)
		if err != nil {
			return fmt.Errorf("ошибка постановки задачи в очередь %s: %w", desc.Seller, err)
		}
	}

	return nil
}

func (s *receiverCoreServiceImpl) setCard(ctx context.Context, in entity.Card) (*entity.Card, error) {
	card, err := s.cardRepo.SelectByVendorID(ctx, in.VendorID)
	if errors.Is(err, ErrObjectNotFound) {
		card, err = s.cardRepo.Insert(ctx, in)
		if err != nil {
			return nil, wrapErr(fmt.Errorf("ошибка при сохранении карточки: %w", err))
		}
	}

	return card, nil
}

func (s *receiverCoreServiceImpl) getCardByVendorID(ctx context.Context, vendorID string) (*entity.Card, error) {
	card, err := s.cardRepo.SelectByVendorID(ctx, vendorID)
	if err != nil {
		return nil, err
	}

	return card, nil
}
