package odincfetcer

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/pkg/alogger"
)

type Card struct {
	VendorID    string `json:"vendor_id"`
	VendorCode  string `json:"vendor_code"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    struct {
		ID    string `json:"ID"`
		Title string `json:"title"`
	} `json:"category"`
	Length  float32 `json:"length"`
	Width   float32 `json:"width"`
	Height  float32 `json:"height"`
	Barcode string  `json:"barcode"`
	Brand   string  `json:"brand"`
	Size    string  `json:"size"`
}

func (odinc *odincAPIclientImp) GetCards(ctx context.Context, desc entity.PackageDescription) ([]entity.Card, error) {
	const methodName = "hs/sender-api/getCardByBarcode"

	queryString := mapToURLValues(desc.Query)

	requestURL := fmt.Sprintf("%s%s?%s", odinc.addr, methodName, queryString.Encode())

	alogger.InfoFromCtx(ctx, "Запрашиваем данные из 1с: %s", desc.Query)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка создания запроса: %w", methodName, err)
	}

	auth := base64.StdEncoding.EncodeToString([]byte(odinc.login + ":" + odinc.password))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Content-Type", "application/json")

	resp, err := odinc.client.Do(req)

	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: ошибка отправки запроса: %w", methodName, err)
	}
	defer resp.Body.Close()

	var odincResponse []Card
	if err := json.NewDecoder(resp.Body).Decode(&odincResponse); err != nil {
		return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", methodName, err)
	}

	var cardsList []entity.Card

	for _, v := range odincResponse {
		brand := entity.Brand{
			Title: v.Brand,
		}

		categories := []*entity.Category{}
		categories = append(categories, &entity.Category{
			Title: v.Category.Title,
		})

		barcodes := []*entity.Barcode{}
		sizes := []*entity.Size{}

		sizes = append(sizes, &entity.Size{
			TechSize: v.Size,
			Title:    v.Size,
		})

		barcodes = append(barcodes, &entity.Barcode{
			Barcode: v.Barcode,
		})

		dimension := entity.Dimension{
			Width:   int(v.Width),
			Height:  int(v.Height),
			Length:  int(v.Length),
			IsVaild: true,
		}

		card := entity.Card{
			VendorID:    "",
			VendorCode:  v.VendorCode,
			Title:       v.Title,
			Description: v.Description,
			Brand:       brand,
			Barcodes:    barcodes,
			Categories:  categories,
			Sizes:       sizes,
			Dimension:   dimension,
			UpdatedAt:   time.Now(),
		}
		cardsList = append(cardsList, card)
	}

	return cardsList, nil
}

func (odinc *odincAPIclientImp) GetStocks(ctx context.Context, desc entity.PackageDescription) ([]entity.StockMeta, error) {
	return nil, nil
}

func (odinc *odincAPIclientImp) GetWarehouses(ctx context.Context) ([]entity.Warehouse, error) {
	return nil, nil
}

func (odinc *odincAPIclientImp) GetOrders(ctx context.Context, desc entity.PackageDescription) ([]entity.Order, error) {
	return nil, nil
}

func (odinc *odincAPIclientImp) GetSales(ctx context.Context, desc entity.PackageDescription) ([]entity.Sale, error) {
	return nil, nil
}

func (odinc *odincAPIclientImp) Ping(ctx context.Context) error {
	return nil
}

func mapToURLValues(data map[string]string) url.Values {
	values := url.Values{}
	for key, value := range data {
		values.Add(key, value)
	}
	return values
}
