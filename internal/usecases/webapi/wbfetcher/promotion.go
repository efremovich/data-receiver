package wbfetcher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
)

func (wb *apiClientImp) GetPromotion(ctx context.Context, desc entity.PackageDescription) ([]entity.Promotion, error) {
	// 1 - Рекламные кампании
	// 2 - Медиакампании
	promotionsList, err := wb.GetPromotionList(ctx, 1)
	if err != nil {
		return nil, err
	}

	return promotionsList, nil
}

func (wb *apiClientImp) GetPromotionList(ctx context.Context, promoType int) ([]entity.Promotion, error) {
	var promotionsList []entity.Promotion

	method := promotionListMethod

	switch promoType {
	case 1:
		method = promotionListMethod
	case 2:
		method = mediaPromotionListMethod
	}

	reqURL := fmt.Sprintf("%s%s", advertAPIURL, method)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка создания запроса: %w", method, err)
	}

	req.Header.Set("Authorization", wb.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := wb.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка отправки запроса: %w", method, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: сервер ответил: %d", method, resp.StatusCode)
	}

	defer resp.Body.Close()

	var response PromotionsList
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", method, err)
	}

	for _, advert := range response.Adverts {
		for _, advertElem := range advert.AdvertList {
			promotion, err := wb.GetPromotionInfo(ctx, advertElem.AdvertID, promoType)
			if err != nil {
				return nil, err
			}

			err = wb.GetPromotionStats(ctx, &promotion)
			if err != nil {
				return nil, err
			}

			promotionsList = append(promotionsList, promotion)
		}
	}

	return promotionsList, nil
}

func (wb *apiClientImp) GetPromotionInfo(ctx context.Context, advertId, promoType int) (entity.Promotion, error) {
	var promotion entity.Promotion

	method := promotionInfoMetod

	switch promoType {
	case 1:
		method = promotionInfoMetod
	case 2:
		method = advertsInfoMetod
	}

	reqURL := fmt.Sprintf("%s%s", advertAPIURL, method)
	filter := []int{
		advertId,
	}

	bodyData, err := json.Marshal(filter)
	if err != nil {
		return promotion, fmt.Errorf("%s: ошибка маршалинга тела запроса: %w", method, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(bodyData))
	if err != nil {
		return entity.Promotion{}, fmt.Errorf("%s: ошибка создания запроса: %w", method, err)
	}

	req.Header.Set("Authorization", wb.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := wb.client.Do(req)
	if err != nil {
		return promotion, fmt.Errorf("%s: ошибка отправки запроса: %w", method, err)
	}

	if resp.StatusCode != http.StatusOK {
		return promotion, fmt.Errorf("%s: сервер ответил: %d", method, resp.StatusCode)
	}

	defer resp.Body.Close()

	var response []AdvertList
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return promotion, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", promotionListMethod, err)
	}

	for _, advert := range response {
		promotion.ExternalID = int64(advert.AdvertID)
		promotion.Name = advert.Name
		promotion.Type = advert.Type
		promotion.DateStart = advert.StartTime
		promotion.DateEnd = advert.EndTime
		promotion.CreateTime = advert.CreateTime
		promotion.ChangeTime = advert.ChangeTime
		promotion.Status = advert.Status
	}

	return promotion, nil
}

func (wb *apiClientImp) GetPromotionStats(ctx context.Context, promotion *entity.Promotion) error {
	type Filter struct {
		ID       int64 `json:"id"`
		Interval struct {
			Begin string `json:"begin"`
			End   string `json:"end"`
		} `json:"interval"`
	}

	// Функция для разбиения периода на интервалы по 31 дню
	splitIntervals := func(start, end time.Time) []struct {
		Begin time.Time
		End   time.Time
	} {
		var intervals []struct {
			Begin time.Time
			End   time.Time
		}

		for current := start; current.Before(end); {
			next := current.AddDate(0, 0, 30) // Добавляем 30 дней, чтобы получить 31-дневный интервал
			if next.After(end) {
				next = end
			}

			intervals = append(intervals, struct {
				Begin time.Time
				End   time.Time
			}{
				Begin: current,
				End:   next,
			})

			current = next.AddDate(0, 0, 1) // Начинаем следующий интервал со следующего дня
		}

		return intervals
	}

	// Получаем все интервалы для запросов
	intervals := splitIntervals(promotion.DateStart, promotion.DateEnd)

	// Собираем статистику из всех запросов
	var allStats []PromotionListDetail

	// Делаем запросы для каждого интервала
	for _, interval := range intervals {
		filter := []Filter{
			{
				ID: promotion.ExternalID,
				Interval: struct {
					Begin string `json:"begin"`
					End   string `json:"end"`
				}{
					Begin: interval.Begin.Format("2006-01-02"),
					End:   interval.End.Format("2006-01-02"),
				},
			},
		}

		reqURL := fmt.Sprintf("%s%s", advertAPIURL, promoFullStatsMethod)

		bodyData, err := json.Marshal(filter)
		if err != nil {
			return fmt.Errorf("%s: ошибка маршалинга тела запроса: %w", promoFullStatsMethod, err)
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(bodyData))
		if err != nil {
			return fmt.Errorf("%s: ошибка создания запроса: %w", promoFullStatsMethod, err)
		}

		req.Header.Set("Authorization", wb.token)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		resp, err := wb.client.Do(req)
		if err != nil {
			return fmt.Errorf("%s: ошибка отправки запроса: %w", promoFullStatsMethod, err)
		}

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("%s: сервер ответил: %d", promoFullStatsMethod, resp.StatusCode)
		}
		defer resp.Body.Close()

		var response []PromotionListDetail
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			return fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", promoFullStatsMethod, err)
		}

		allStats = append(allStats, response...)

		time.Sleep(61 * time.Second)
	}

	// Обрабатываем все полученные данные
	for _, elem := range allStats {
		promotion.Views += elem.Views
		promotion.Clicks += elem.Clicks
		promotion.CTR += elem.Ctr
		promotion.CPC += elem.Cpc
		promotion.Spent += elem.Sum
		promotion.Orders += elem.Orders
		promotion.CR += elem.Cr
		promotion.SHKs += elem.Shks
		promotion.OrderAmount += elem.SumPrice

		for _, days := range elem.Days {
			for _, apps := range days.Apps {
				for _, cardPromoDetail := range apps.Nm {
					promoStats := entity.PromotionStats{}
					promoStats.Views = cardPromoDetail.Views
					promoStats.Clicks = cardPromoDetail.Clicks
					promoStats.CTR = cardPromoDetail.Ctr
					promoStats.CPC = cardPromoDetail.Cpc
					promoStats.Spent = cardPromoDetail.Sum
					promoStats.Orders = cardPromoDetail.Orders
					promoStats.CR = cardPromoDetail.Cr
					promoStats.SHKs = cardPromoDetail.Shks
					promoStats.OrderAmount = cardPromoDetail.SumPrice
					promoStats.PromotionID = promotion.ID
					promoStats.CardExternalID = int64(cardPromoDetail.NmID)
					promoStats.Date = days.Date
					promoStats.AppType = apps.AppType

					promotion.PromotionStats = append(promotion.PromotionStats, promoStats)
				}
			}
		}
	}

	return nil
}
