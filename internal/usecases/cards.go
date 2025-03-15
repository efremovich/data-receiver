package usecases

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/webapi"
	"github.com/efremovich/data-receiver/pkg/alogger"
)

func (s *receiverCoreServiceImpl) ReceiveCards(ctx context.Context, desc entity.PackageDescription) error {
	clients := s.apiFetcher[desc.Seller]

	g, gCtx := errgroup.WithContext(ctx)

	for _, c := range clients {
		client := c

		g.Go(func() error {
			return s.receiveAndSaveCard(gCtx, client, desc)
		})
	}

	if err := g.Wait(); err != nil {
		if errors.Is(err, context.Canceled) {
			alogger.WarnFromCtx(ctx, "Операция была отменена: %v", err)
			return nil
		}
		return fmt.Errorf("ошибка при обработке клиентов: %w", err)
	}

	return nil
}

func (s *receiverCoreServiceImpl) receiveAndSaveCard(ctx context.Context, client webapi.ExtAPIFetcher, desc entity.PackageDescription) error {
	cards, err := client.GetCards(ctx, desc)
	if err != nil {
		return fmt.Errorf("ошибка получение данные из внешнего источника %s, %w", desc.Seller, err)
	}

	// Seller
	seller, err := s.getSeller(ctx, client.GetMarketPlace())
	if err != nil {
		return err
	}

	alogger.InfoFromCtx(ctx, "Начали загрузку карточек товара %d от %s", len(cards), seller.Title)

	for _, in := range cards {
		s.metricsCollector.IncServiceDocsTaskCounter()

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

	alogger.InfoFromCtx(ctx, "Задание загрузка карточек товаров успешно завешена количество: %d, маркетплейс %s", len(cards), seller.Title)

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

func (s *receiverCoreServiceImpl) getCardByID(ctx context.Context, cardID int64) (*entity.Card, error) {
	card, err := s.cardRepo.SelectByID(ctx, cardID)
	if err != nil {
		return nil, err
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
