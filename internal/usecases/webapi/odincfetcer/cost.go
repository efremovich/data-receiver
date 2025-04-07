package odincfetcer

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/pkg/alogger"
)

const getCostsMetod = "hs/sender-api/cost_price"
const LIMIT = 1000

func (odinc *apiClientImp) GetCosts(ctx context.Context, desc entity.PackageDescription) ([]entity.Cost, error) {
	queryString := url.Values{}
	queryString.Set("date", desc.UpdatedAt.Format("2006-01-02"))

	cursor := ""

	header := http.Header{}
	auth := base64.StdEncoding.EncodeToString([]byte(odinc.login + ":" + odinc.password))
	header.Set("Authorization", "Basic "+auth)
	header.Set("Content-Type", "application/json")

	var costList []entity.Cost

	run := true
	for run {
		if cursor != "" {
			queryString.Set("cursor", cursor)
		}

		requestURL := fmt.Sprintf("%s%s?%s", marketPlaceAPIURL, getCostsMetod, queryString.Encode())

		alogger.InfoFromCtx(ctx, "Запрашиваем данные из 1с: %+v", queryString)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
		if err != nil {
			return nil, fmt.Errorf("%s: ошибка создания запроса: %w", getCostsMetod, err)
		}

		req.Header = header

		resp, err := odinc.client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("%s: ошибка отправки запроса: %w", getCostsMetod, err)
		}

		defer resp.Body.Close()

		response := []Cost{}

		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return nil, fmt.Errorf("%s: ошибка отправки запроса: %w", getCostsMetod, err)
		}

		if len(response) < LIMIT {
			run = false
		}

		for _, elem := range response {
			costPrice, err := strconv.ParseFloat(elem.CostPrice, 32)
			if err != nil {
				return nil, fmt.Errorf("ошибка конвертации себестоимости в float32: %w", err)
			}

			cost := entity.Cost{}
			cost.Amount = costPrice
			cost.ExternalID = elem.VendorCode
			cost.CreatedAt = desc.UpdatedAt

			cursor = strings.TrimSpace(elem.VendorCode)

			costList = append(costList, cost)
		}
	}

	return costList, nil
}
