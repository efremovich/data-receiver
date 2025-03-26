package usecases

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
)

func (s *receiverCoreServiceImpl) OfferFeed(ctx context.Context) ([]byte, error) {
	offers, err := s.offerfeedrepo.GetOffers(ctx)
	if err != nil {
		return nil, err
	}

	yml := entity.YMLCatalog{
		Date: time.Now().Format("2006-01-02 15:04:05"),
		Shop: entity.Shop{
			Name:       "Название магазина",
			Company:    "Название компании",
			URL:        "Ссылка на сайт",
			Categories: nil,
			Offers:     offers,
		},
	}

	output, err := xml.MarshalIndent(yml, "", "  ")
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (s *receiverCoreServiceImpl) StockFeed(ctx context.Context) ([]byte, error) {
	inventory, err := s.offerfeedrepo.GetStocks(ctx)
	if err != nil {
		return nil, err
	}

	yml := entity.YMLCatalog{
		Date: time.Now().Format("2006-01-02 15:04:05"),
		Shop: entity.Shop{
			Name:      "Название магазина",
			Company:   "Название компании",
			URL:       "Ссылка на сайт",
			Inventory: inventory,
		},
	}

	output, err := xml.MarshalIndent(yml, "", "  ")
	if err != nil {
		return nil, err
	}

	return output, nil
}

func (s *receiverCoreServiceImpl) VkCardsFeed(ctx context.Context, params entity.VkCardsFeedParams) ([]byte, error) {
	vkCards, err := s.offerfeedrepo.GetCardsVkFeed(ctx, params)
	if err != nil {
		return nil, err
	}

	cursor := ""

	total := len(vkCards)
	if total > 0 {
		cursor = vkCards[total-1].VendorCode
	}

	feedResponse := entity.ResponseVKCard{
		Cards:  vkCards,
		Total:  total,
		Cursor: cursor,
	}

	output, err := json.Marshal(feedResponse)
	if err != nil {
		return nil, err
	}

	return output, nil
}
