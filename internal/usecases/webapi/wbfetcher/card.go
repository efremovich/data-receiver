package wbfetcher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
)

var reVendorCode = regexp.MustCompile(`\d{2}-\d{7,8}`)

type WbResponse struct {
	Cards  []Cards `json:"cards"`
	Cursor Cursor  `json:"cursor"`
}
type Photos struct {
	Big      string `json:"big"`
	C246X328 string `json:"c246x328"`
	C516X688 string `json:"c516x688"`
	Square   string `json:"square"`
	Tm       string `json:"tm"`
}
type Dimensions struct {
	Length  int  `json:"length"`
	Width   int  `json:"width"`
	Height  int  `json:"height"`
	IsValid bool `json:"isValid"`
}
type Characteristics struct {
	ID    int         `json:"id"`
	Name  string      `json:"name"`
	Value interface{} `json:"value,omitempty"`
}
type Sizes struct {
	ChrtID   int      `json:"chrtID"`
	TechSize string   `json:"techSize"`
	WbSize   string   `json:"wbSize"`
	Skus     []string `json:"skus"`
}
type Tags struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}
type Cards struct {
	NmID            int               `json:"nmID"`
	ImtID           int               `json:"imtID"`
	NmUUID          string            `json:"nmUUID"`
	SubjectID       int               `json:"subjectID"`
	SubjectName     string            `json:"subjectName"`
	VendorCode      string            `json:"vendorCode"`
	Brand           string            `json:"brand"`
	Title           string            `json:"title"`
	Description     string            `json:"description"`
	Photos          []Photos          `json:"photos"`
	Video           string            `json:"video"`
	Dimensions      Dimensions        `json:"dimensions"`
	Characteristics []Characteristics `json:"characteristics"`
	Sizes           []Sizes           `json:"sizes"`
	Tags            []Tags            `json:"tags"`
	CreatedAt       time.Time         `json:"createdAt"`
	UpdatedAt       time.Time         `json:"updatedAt"`
}

type Cursor struct {
	NmID      int        `json:"external_id,omitempty"`
	Total     int        `json:"total,omitempty"`
	Limit     int        `json:"limit,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

type Filter struct {
	WithPhoto int `json:"withPhoto"`
}

type Sort struct {
	Ascending bool `json:"ascending"` // true - asc sort, false - desc sort
}

type Settings struct {
	Sort   Sort   `json:"sort"`
	Cursor Cursor `json:"cursor"`
	Filter Filter `json:"filter"`
}

type Setting struct {
	Setting Settings `json:"settings"`
}

func (wb *wbAPIclientImp) GetCards(ctx context.Context, desc entity.PackageDescription) ([]entity.Card, error) {
	const methodName = "/content/v2/get/cards/list?locale=ru"

	lastID, err := strconv.Atoi(desc.Cursor)
	if err != nil {
		return nil, fmt.Errorf("ошибка конвертации lastID: %s,  %w", desc.Cursor, err)
	}

	requestSettings := Settings{
		Sort:   Sort{Ascending: true},
		Filter: Filter{WithPhoto: -1},
		Cursor: Cursor{Limit: desc.Limit, NmID: lastID, UpdatedAt: &desc.UpdatedAt},
	}

	requestData, err := json.Marshal(Setting{Setting: requestSettings})
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка маршалинга тела запроса: %w", methodName, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s%s", wb.addrContent, methodName), bytes.NewReader(requestData))
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка создания запроса: %w", methodName, err)
	}

	req.Header.Set("Authorization", wb.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")

	resp, err := wb.client.Do(req)

	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: ошибка отправки запроса: %w", methodName, err)
	}
	defer resp.Body.Close()

	var wbResponse WbResponse
	if err := json.NewDecoder(resp.Body).Decode(&wbResponse); err != nil {
		return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", methodName, err)
	}

	var cardsList []entity.Card

	for _, v := range wbResponse.Cards {
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

func (wb *wbAPIclientImp) Ping(_ context.Context) error {
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
