package wbfetcher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/rand"
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
			if limitRemaning := resp.Header.Get("X-Ratelimit-Reset"); limitRemaning == "0" {
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

	maxRetries := 5               // Увеличили количество попыток
	retryDelay := 1 * time.Second // Начальная задержка

	const maxDelay = 60 * time.Second // Максимальная задержка

	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		// Создаем новый контекст с таймаутом для каждого запроса
		reqCtx, cancel := context.WithTimeout(ctx, 120*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(reqCtx, http.MethodPost, reqURL, bytes.NewReader(bodyData))
		if err != nil {
			return fmt.Errorf("%s: ошибка создания запроса: %w", promoFullStatsMethod, err)
		}

		req.Header.Set("Authorization", wb.token)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		// Добавляем логирование перед отправкой запроса
		alogger.DebugFromCtx(ctx, "Отправка запроса (попытка %d/%d)", attempt, maxRetries)

		resp, err := wb.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("%s: ошибка отправки запроса: %w", promoFullStatsMethod, err)

			if attempt < maxRetries {
				// Экспоненциальный backoff с джиттером
				retryDelay = exponentialBackoff(attempt, maxDelay)
				alogger.DebugFromCtx(ctx, "Ошибка запроса, повтор через %v: %v", retryDelay, err)
				time.Sleep(retryDelay)
			}

			continue
		}

		// Обработка rate limiting
		if resp.StatusCode == http.StatusTooManyRequests {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			// Пытаемся получить время ожидания из заголовков
			retryAfter := parseRetryAfter(resp.Header)
			if retryAfter == 0 {
				retryAfter = exponentialBackoff(attempt, maxDelay)
			}

			lastErr = fmt.Errorf("%s: сервер ответил: %d %s (Too Many Requests)", promoFullStatsMethod, resp.StatusCode, string(body))

			if attempt < maxRetries {
				alogger.DebugFromCtx(ctx, "Превышен лимит запросов, повтор через %v", retryAfter)
				time.Sleep(retryAfter)
			}

			continue
		}

		// Обработка других ошибок HTTP
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			lastErr = fmt.Errorf("%s: сервер ответил: %d %s", promoFullStatsMethod, resp.StatusCode, string(body))

			if attempt < maxRetries && shouldRetry(resp.StatusCode) {
				retryAfter := exponentialBackoff(attempt, maxDelay)
				alogger.DebugFromCtx(ctx, "Ошибка %d, повтор через %v", resp.StatusCode, retryAfter)
				time.Sleep(retryAfter)

				continue
			}

			if resp.StatusCode >= 400 && resp.StatusCode < 500 && resp.StatusCode != http.StatusTooManyRequests {
				return fmt.Errorf("%s: сервер ответил: %d %s %w", promoFullStatsMethod, resp.StatusCode, string(body), entity.ErrPermanent)
			}

			continue
		}

		// Проверяем оставшийся лимит запросов
		if remaining := resp.Header.Get("X-Ratelimit-Remaining"); remaining == "0" {
			if reset := resp.Header.Get("X-Ratelimit-Reset"); reset != "" {
				if resetSec, err := strconv.Atoi(reset); err == nil {
					alogger.DebugFromCtx(ctx, "Достигнут лимит запросов, ожидание %d секунд", resetSec)
					time.Sleep(time.Duration(resetSec) * time.Second)
				}
			}
		}

		var response []PromotionListDetail
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			resp.Body.Close()

			return fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", promoFullStatsMethod, err)
		}

		resp.Body.Close()

		alogger.InfoFromCtx(ctx, "обработка данных по рекламным компаниям %d", len(response))

		// Выносим обработку ответа в отдельную функцию для улучшения читаемости
		processPromotionResponse(ctx, promotions, response)

		return nil
	}

	return lastErr
}

// Вспомогательные функции

// exponentialBackoff вычисляет время задержки с экспоненциальным ростом и джиттером.
func exponentialBackoff(attempt int, maxDelay time.Duration) time.Duration {
	delay := time.Duration(math.Pow(2, float64(attempt))) * time.Second
	if delay > maxDelay {
		delay = maxDelay
	}
	// Добавляем случайный джиттер (10% от delay)
	jitter := time.Duration(rand.Int63n(int64(delay / 10)))
	return delay + jitter
}

// parseRetryAfter парсит заголовки для определения времени ожидания.
func parseRetryAfter(headers http.Header) time.Duration {
	// Сначала проверяем X-Ratelimit-Retry
	if retry := headers.Get("X-Ratelimit-Retry"); retry != "" {
		if sec, err := strconv.Atoi(retry); err == nil {
			return time.Duration(sec) * time.Second
		}
	}

	// Затем проверяем X-Ratelimit-Reset
	if reset := headers.Get("X-Ratelimit-Reset"); reset != "" {
		if sec, err := strconv.Atoi(reset); err == nil {
			return time.Duration(sec) * time.Second
		}
	}

	// Затем стандартный Retry-After
	if retryAfter := headers.Get("Retry-After"); retryAfter != "" {
		if sec, err := strconv.Atoi(retryAfter); err == nil {
			return time.Duration(sec) * time.Second
		}
	}

	return 0
}

// shouldRetry определяет, стоит ли повторять запрос при данном коде состояния.
func shouldRetry(statusCode int) bool {
	return statusCode >= 500 ||
		statusCode == http.StatusTooManyRequests ||
		statusCode == http.StatusRequestTimeout ||
		statusCode == http.StatusConflict // 409 иногда временная ошибка
}

// processPromotionResponse выносит логику обработки ответа в отдельную функцию
func processPromotionResponse(ctx context.Context, promotions map[int64]*entity.Promotion, response []PromotionListDetail) {
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
					promoStats := entity.PromotionStats{
						Views:          cardPromoDetail.Views,
						Clicks:         cardPromoDetail.Clicks,
						CTR:            cardPromoDetail.Ctr,
						CPC:            cardPromoDetail.Cpc,
						Spent:          cardPromoDetail.Sum,
						Orders:         cardPromoDetail.Orders,
						CR:             cardPromoDetail.Cr,
						SHKs:           cardPromoDetail.Shks,
						OrderAmount:    cardPromoDetail.SumPrice,
						PromotionID:    promotion.ID,
						CardExternalID: int64(cardPromoDetail.NmID),
						Date:           days.Date,
						AppType:        apps.AppType,
					}
					promotion.PromotionStats = append(promotion.PromotionStats, promoStats)
				}
			}
		}
	}
}

// Разбивает интервал на подынтервалы по maxDays дней.
func splitInterval(promoFilter PromotionFilter, maxDays int) []PromotionFilter {
	begin, _ := time.Parse("2006-01-02", promoFilter.Interval.Begin)
	end, _ := time.Parse("2006-01-02", promoFilter.Interval.End)
	days := int(end.Sub(begin).Hours() / 24)

	if days <= maxDays {
		return []PromotionFilter{promoFilter}
	}

	var chunks []PromotionFilter
	currentBegin := begin

	for {
		currentEnd := currentBegin.AddDate(0, 0, maxDays-1)
		if currentEnd.After(end) {
			currentEnd = end
		}

		chunks = append(chunks, PromotionFilter{
			ID: promoFilter.ID,
			Interval: struct {
				Begin string `json:"begin"`
				End   string `json:"end"`
			}{
				Begin: currentBegin.Format("2006-01-02"),
				End:   currentEnd.Format("2006-01-02"),
			},
		})

		if currentEnd.Equal(end) {
			break
		}
		currentBegin = currentEnd.AddDate(0, 0, 1)
	}

	return chunks
}

// Делит срез на чанки, максимально заполняя каждый 100 уникальными элементами
func chunkPromotionFilters(filters []PromotionFilter, chunkSize, maxDays int) [][]PromotionFilter {
	// 1. Сначала группируем все фильтры по их ID
	filtersByID := make(map[int64][]PromotionFilter)

	for _, filter := range filters {
		begin, err1 := time.Parse("2006-01-02", filter.Interval.Begin)
		end, err2 := time.Parse("2006-01-02", filter.Interval.End)

		if err1 != nil || err2 != nil {
			continue
		}

		days := int(end.Sub(begin).Hours() / 24)
		if days > maxDays {
			chunks := splitInterval(filter, maxDays)
			filtersByID[filter.ID] = append(filtersByID[filter.ID], chunks...)
		} else {
			filtersByID[filter.ID] = append(filtersByID[filter.ID], filter)
		}
	}

	// 2. Создаем очередь для каждого ID с его фильтрами
	idQueue := make([]int64, 0, len(filtersByID))
	filtersQueue := make(map[int64][]PromotionFilter)

	for id, flts := range filtersByID {
		idQueue = append(idQueue, id)
		filtersQueue[id] = flts
	}

	// 3. Распределяем фильтры по чанкам, максимизируя заполнение
	var result [][]PromotionFilter
	currentChunk := make([]PromotionFilter, 0, chunkSize)

	for len(idQueue) > 0 {
		// Берем следующий ID из очереди
		id := idQueue[0]
		idQueue = idQueue[1:]

		// Добавляем первый доступный фильтр для этого ID
		if len(filtersQueue[id]) > 0 {
			currentChunk = append(currentChunk, filtersQueue[id][0])
			filtersQueue[id] = filtersQueue[id][1:]

			// Если остались еще фильтры для этого ID, возвращаем в конец очереди
			if len(filtersQueue[id]) > 0 {
				idQueue = append(idQueue, id)
			}
		}

		// Если чанк заполнен или больше нет элементов, добавляем в результат
		if len(currentChunk) >= chunkSize || len(idQueue) == 0 {
			if len(currentChunk) > 0 {
				result = append(result, currentChunk)
				currentChunk = make([]PromotionFilter, 0, chunkSize)
			}
		}
	}

	return result
}
