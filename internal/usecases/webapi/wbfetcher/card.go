package wbfetcher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/efremovich/data-receiver/internal/entity"
)

var reVendorCode = regexp.MustCompile(`\d{2}-\d{5,8}`)

func (wb *apiClientImp) GetCards(ctx context.Context, _ entity.PackageDescription) ([]entity.Card, error) {
	var cardsList []entity.Card

	cards, err := wb.getCardsFromWB(ctx)
	if err != nil {
		return cardsList, err
	}

	for _, v := range cards {
		brand := entity.Brand{
			Title: v.Brand,
		}

		characteristics := []*entity.CardCharacteristic{}

		for _, c := range v.Characteristics {
			char := entity.CardCharacteristic{
				Title: c.Name,
				Value: convertInterfaceToStringSlice(c.Value),
			}
			characteristics = append(characteristics, &char)
		}

		categories := []*entity.Category{}
		categories = append(categories, &entity.Category{
			Title: v.SubjectName,
		})

		barcodes := []*entity.Barcode{}
		sizes := []*entity.Size{}

		for _, sz := range v.Sizes {
			sizes = append(sizes, &entity.Size{
				TechSize:   sz.TechSize,
				Title:      sz.WbSize,
				ExternalID: int64(sz.ChrtID),
			})

			for _, b := range sz.Skus {
				barcodes = append(barcodes, &entity.Barcode{
					ExternalID: int64(sz.ChrtID),
					Barcode:    b,
				})
			}
		}

		mediaFile := []*entity.MediaFile{}

		for _, mf := range v.Photos {
			mediaFile = append(mediaFile, &entity.MediaFile{
				Link:   mf.Big,
				TypeID: 1,
			})
		}

		if v.Video != "" {
			mediaFile = append(mediaFile, &entity.MediaFile{
				Link:   v.Video,
				TypeID: 2,
			})
		}

		dimension := entity.Dimension{
			Width:   v.Dimensions.Width,
			Height:  v.Dimensions.Height,
			Length:  v.Dimensions.Length,
			IsVaild: v.Dimensions.IsValid,
		}

		vendorCode := v.VendorCode
		if reVendorCode.MatchString(v.VendorCode) {
			vendorCode = reVendorCode.FindString(v.VendorCode)
		}

		card := entity.Card{
			ExternalID:      int64(v.NmID),
			VendorID:        vendorCode,
			VendorCode:      vendorCode,
			Title:           v.Title,
			Description:     v.Description,
			Brand:           brand,
			Characteristics: characteristics,
			Barcodes:        barcodes,
			Categories:      categories,
			Sizes:           sizes,
			MediaFile:       mediaFile,
			Dimension:       dimension,
			UpdatedAt:       v.UpdatedAt,
		}
		cardsList = append(cardsList, card)
	}

	return cardsList, nil
}

func (wb *apiClientImp) getCardsFromWB(ctx context.Context) ([]*Cards, error) {
	requestSettings := Settings{
		Sort:   Sort{Ascending: true},
		Filter: Filter{WithPhoto: -1},
		Cursor: Cursor{Limit: cardRequestLimit},
	}

	cards := []*Cards{}

	run := true

	for run {
		requestData, err := json.Marshal(Setting{Setting: requestSettings})
		if err != nil {
			return nil, fmt.Errorf("%s: ошибка маршалинга тела запроса: %w", cardListMethod, err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s%s", contentAPIURL, cardListMethod), bytes.NewReader(requestData))
		if err != nil {
			return nil, fmt.Errorf("%s: ошибка создания запроса: %w", cardListMethod, err)
		}

		req.Header.Set("Authorization", wb.token)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("accept", "application/json")

		resp, err := wb.client.Do(req)

		if err != nil || resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("%s: ошибка отправки запроса: %w", cardListMethod, err)
		}
		defer resp.Body.Close()

		var wbResponse WbResponse
		if err := json.NewDecoder(resp.Body).Decode(&wbResponse); err != nil {
			return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", cardListMethod, err)
		}

		for _, card := range wbResponse.Cards {
			cards = append(cards, &card)
		}

		if wbResponse.Cursor.Total < cardRequestLimit {
			run = false
		}

		requestSettings.Cursor.NmID = wbResponse.Cursor.NmID
		requestSettings.Cursor.UpdatedAt = wbResponse.Cursor.UpdatedAt
	}

	return cards, nil
}

func (wb *apiClientImp) Ping(_ context.Context) error {
	return nil
}

func convertInterfaceToStringSlice(input interface{}) []string {
	switch v := input.(type) {
	case int:
		return []string{strconv.Itoa(v)}
	case string:
		return []string{v}
	case []string:
		return v
	case []interface{}:
		var output []string
		for _, vv := range v {
			output = append(output, convertInterfaceToStringSlice(vv)...)
		}

		return output
	}

	return nil
}
