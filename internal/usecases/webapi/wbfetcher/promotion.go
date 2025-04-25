package wbfetcher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/pkg/alogger"
)

func (wb *apiClientImp) GetPromotion(ctx context.Context, desc entity.PackageDescription) ([]*entity.Promotion, error) {
	// 1 - Рекламные кампании
	// 2 - Медиакампании
	promotionsList, err := wb.GetPromotionList(ctx, 1)
	if err != nil {
		return nil, err
	}

	return promotionsList, nil
}

func (wb *apiClientImp) GetPromotionList(ctx context.Context, promoType int) ([]*entity.Promotion, error) {
	var promotionsList []*entity.Promotion

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
		return nil, fmt.Errorf("%s: сервер ответил: %d %w", method, resp.StatusCode, entity.ErrPermanent)
	}

	var response PromotionsList
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", method, err)
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка закрытия тела ответа: %w", method, err)
	}

	ticker := time.NewTicker(200 * time.Millisecond) // 5 раз в секунду (1000ms/5 = 200ms)
	defer ticker.Stop()

	alogger.InfoFromCtx(ctx, "Начали загрузку данных по рекламной компании %d", response.All)

	for i := range response.Adverts {
		for z := range response.Adverts[i].AdvertList {
			<-ticker.C

			promotion, err := wb.GetPromotionInfo(ctx, response.Adverts[i].AdvertList[z].AdvertID, promoType)
			if err != nil {
				return nil, err
			}

			promotionsList = append(promotionsList, promotion)
		}
	}

	filters := []PromotionFilter{}

	for _, promotion := range promotionsList {
		if promotion.Status == 11 { // Компания завершена
			continue // Не обрабатываем, так как нет данных
		}

		if promotion.DateEnd.After(time.Now()) {
			promotion.DateEnd = time.Now()
		}

		filter := PromotionFilter{
			ID: promotion.ExternalID,
			Interval: struct {
				Begin string `json:"begin"`
				End   string `json:"end"`
			}{
				Begin: promotion.DateStart.Format("2006-01-02"),
				End:   promotion.DateEnd.Format("2006-01-02"),
			},
		}

		filters = append(filters, filter)
	}

	chunkFilter := chunkPromotionFilters(filters, 100, 31)

	promoMap := make(map[int64]*entity.Promotion, len(promotionsList))
	for _, promotion := range promotionsList {
		promoMap[promotion.ExternalID] = promotion
	}

	for i := range chunkFilter {
		alogger.InfoFromCtx(ctx, "Получение дальной информации по компании %d из %d", i, len(chunkFilter))

		err = wb.GetPromotionStats(ctx, promoMap, chunkFilter[i])
		if err != nil {
			return nil, err
		}
	}

	return promotionsList, nil
}

func (wb *apiClientImp) GetPromotionInfo(ctx context.Context, advertId, promoType int) (*entity.Promotion, error) {
	promotion := new(entity.Promotion)
	method := promotionInfoMetod

	switch promoType {
	case 1:
		method = promotionInfoMetod
	case 2:
		method = advertsInfoMetod
	}

	reqURL := fmt.Sprintf("%s%s", advertAPIURL, method)
	filter := []int{advertId}

	bodyData, err := json.Marshal(filter)
	if err != nil {
		return promotion, fmt.Errorf("%s: ошибка маршалинга тела запроса: %w", method, err)
	}

	maxRetries := 3
	retryDelay := time.Second * 60

	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(bodyData))
		if err != nil {
			return nil, fmt.Errorf("%s: ошибка создания запроса: %w", method, err)
		}

		req.Header.Set("Authorization", wb.token)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		resp, err := wb.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("%s: ошибка отправки запроса: %w", method, err)

			if attempt < maxRetries {
				time.Sleep(retryDelay)
			}

			continue
		}

		// Обработка статуса 429
		if resp.StatusCode == http.StatusTooManyRequests {
			err := resp.Body.Close()
			if err != nil {
				return promotion, fmt.Errorf("%s: ошибка закрытия тела ответа: %w", method, err)
			}
			// Пытаемся получить время ожидания из заголовка Retry-After
			retryAfter := resp.Header.Get("X-Ratelimit-Reset")
			if retryAfter != "" {
				if seconds, err := strconv.Atoi(retryAfter); err == nil {
					retryDelay = time.Second * time.Duration(seconds)
				}
			}

			lastErr = fmt.Errorf("%s: сервер ответил: %d (Too Many Requests)", method, resp.StatusCode)

			if attempt < maxRetries {
				time.Sleep(retryDelay)
			}

			continue
		}

		// Обработка других ошибок HTTP
		if resp.StatusCode != http.StatusOK {
			err := resp.Body.Close()
			if err != nil {
				return promotion, fmt.Errorf("%s: ошибка закрытия тела ответа: %w", method, err)
			}

			return promotion, fmt.Errorf("%s: сервер ответил: %d %w", method, resp.StatusCode, entity.ErrPermanent)
		}

		if resp.StatusCode == http.StatusOK {
			if limitRemaning := resp.Header.Get("X-Ratelimit-Remaining"); limitRemaning == "0" {
				time.Sleep(retryDelay)
			}
		}

		var response []AdvertList
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			err := resp.Body.Close()
			if err != nil {
				return promotion, fmt.Errorf("%s: ошибка закрытия тела ответа: %w", method, err)
			}

			return promotion, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", promotionListMethod, err)
		}

		err = resp.Body.Close()
		if err != nil {
			return promotion, fmt.Errorf("%s: ошибка закрытия тела ответа: %w", promotionListMethod, err)
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

	return promotion, lastErr
}

func (wb *apiClientImp) GetPromotionStats(ctx context.Context, promotions map[int64]*entity.Promotion, chunks []PromotionFilter) error {
	reqURL := fmt.Sprintf("%s%s", advertAPIURL, promoFullStatsMethod)

	bodyData, err := json.Marshal(chunks)
	if err != nil {
		return fmt.Errorf("%s: ошибка маршалинга тела запроса: %w", promoFullStatsMethod, err)
	}

	maxRetries := 3
	retryDelay := time.Second * 60

	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(bodyData))
		if err != nil {
			return fmt.Errorf("%s: ошибка создания запроса: %w", promoFullStatsMethod, err)
		}

		req.Header.Set("Authorization", wb.token)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		resp, err := wb.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("%s: ошибка отправки запроса: %w", promoFullStatsMethod, err)

			if attempt < maxRetries {
				time.Sleep(retryDelay)
			}

			continue
		}

		// Обработка статуса 429
		if resp.StatusCode == http.StatusTooManyRequests {
			body, _ := io.ReadAll(resp.Body)

			err := resp.Body.Close()
			if err != nil {
				return fmt.Errorf("%s: ошибка закрытия тела ответа: %w", promoFullStatsMethod, err)
			}

			// Пытаемся получить время ожидания из заголовка X-Ratelimit-Retry
			if retryAfter := resp.Header.Get("X-Ratelimit-Reset"); retryAfter != "" {
				if seconds, err := strconv.Atoi(retryAfter); err == nil {
					retryDelay = time.Second * time.Duration(seconds)
				}
			}

			lastErr = fmt.Errorf("%s: сервер ответил: %d %s (Too Many Requests)", promoFullStatsMethod, resp.StatusCode, string(body))

			if attempt < maxRetries {
				time.Sleep(retryDelay)
			}

			continue
		}

		// Обработка других ошибок HTTP
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)

			err := resp.Body.Close()
			if err != nil {
				return fmt.Errorf("%s: ошибка закрытия тела ответа: %w", promoFullStatsMethod, err)
			}

			lastErr = fmt.Errorf("%s: сервер ответил: %d %s", promoFullStatsMethod, resp.StatusCode, string(body))

			if attempt < maxRetries && resp.StatusCode >= 500 {
				// Повторяем только для 5xx ошибок
				time.Sleep(retryDelay)

				continue
			}

			return fmt.Errorf("%s: сервер ответил: %d %s %w", promoFullStatsMethod, resp.StatusCode, string(body), entity.ErrPermanent)
		}

		if resp.StatusCode == http.StatusOK {
			if limitRemaning := resp.Header.Get("X-Ratelimit-Remaining"); limitRemaning == "0" {
				time.Sleep(retryDelay)
			}
		}

		var response []PromotionListDetail
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			err := resp.Body.Close()
			if err != nil {
				return fmt.Errorf("%s: ошибка закрытия тела ответа: %w", promoFullStatsMethod, err)
			}

			return fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", promoFullStatsMethod, err)
		}

		err = resp.Body.Close()
		if err != nil {
			return fmt.Errorf("%s: ошибка закрытия тела ответа: %w", promoFullStatsMethod, err)
		}

		alogger.InfoFromCtx(ctx, "обработка данных по рекламным компаниям %d", len(response))
		// Обрабатываем все полученные данные
		for _, elem := range response {
			promotion, ok := promotions[int64(elem.AdvertID)]
			if !ok {
				alogger.WarnFromCtx(ctx, "Промо %d не нашлось в карте", elem.AdvertID)

				continue
			}

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

	return lastErr
}

// Разбивает интервал на подынтервалы по maxDays дней.
func splitInterval(promoFilter PromotionFilter, maxDays int) []PromotionFilter {
	begin, _ := time.Parse("2006-01-02", promoFilter.Interval.Begin)
	end, _ := time.Parse("2006-01-02", promoFilter.Interval.End)
	days := int(end.Sub(begin).Hours() / 24)

	if days <= maxDays {
		return []PromotionFilter{promoFilter} // Не нужно разбивать
	}

	var chunks []PromotionFilter

	currentBegin := begin

	for {
		currentEnd := currentBegin.AddDate(0, 0, maxDays-1) // 31 день = 30 дней разницы
		if currentEnd.After(end) {
			currentEnd = end
		}

		chunk := PromotionFilter{
			ID: promoFilter.ID,
			Interval: struct {
				Begin string `json:"begin"`
				End   string `json:"end"`
			}{
				Begin: currentBegin.Format("2006-01-02"),
				End:   currentEnd.Format("2006-01-02"),
			},
		}
		chunks = append(chunks, chunk)

		if currentEnd.Equal(end) {
			break
		}

		currentBegin = currentEnd.AddDate(0, 0, 1) // Следующий день
	}

	return chunks
}

// Делит срез на чанки, группируя по 100 элементов, но без дубликатов ID в одном чанке.
func chunkPromotionFilters(filters []PromotionFilter, chunkSize, maxDays int) [][]PromotionFilter {
	var allChunks []PromotionFilter

	// 1. Разбиваем все интервалы > maxDays на подынтервалы.
	for _, promoFilter := range filters {
		begin, err1 := time.Parse("2006-01-02", promoFilter.Interval.Begin)
		end, err2 := time.Parse("2006-01-02", promoFilter.Interval.End)

		if err1 != nil || err2 != nil {
			continue // Пропускаем некорректные данные
		}

		days := int(end.Sub(begin).Hours() / 24)
		if days > maxDays {
			chunks := splitInterval(promoFilter, maxDays)
			allChunks = append(allChunks, chunks...)
		} else {
			allChunks = append(allChunks, promoFilter)
		}
	}

	// 2. Группируем в чанки, избегая дубликатов ID внутри одного чанка.
	var resultChunks [][]PromotionFilter

	currentChunk := make([]PromotionFilter, 0, chunkSize)
	idsInCurrentChunk := make(map[int64]bool) // Для проверки дубликатов

	for _, promoFilter := range allChunks {
		// Если ID уже есть в текущем чанке или чанк заполнен, создаём новый.
		if idsInCurrentChunk[promoFilter.ID] || len(currentChunk) >= chunkSize {
			resultChunks = append(resultChunks, currentChunk)
			currentChunk = make([]PromotionFilter, 0, chunkSize)
			idsInCurrentChunk = make(map[int64]bool)
		}

		currentChunk = append(currentChunk, promoFilter)
		idsInCurrentChunk[promoFilter.ID] = true
	}

	// Добавляем последний чанк, если он не пустой.
	if len(currentChunk) > 0 {
		resultChunks = append(resultChunks, currentChunk)
	}

	return resultChunks
}
