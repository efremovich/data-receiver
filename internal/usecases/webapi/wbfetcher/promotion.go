package wbfetcher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

	for i := range 2 { // response.Adverts {
		for z := range 1 { // response.Adverts[i].AdvertList {
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

	chunkFilter := chunkPromotionFilters(filters, 10, 30)

	promoMap := make(map[int64]*entity.Promotion, len(promotionsList))
	for _, promotion := range promotionsList {
		promoMap[promotion.ExternalID] = promotion
	}

	for i := range chunkFilter {
		err = wb.GetPromotionStats(ctx, promoMap, chunkFilter[i])
		if err != nil {
			return nil, err
		}

		ticker := time.NewTicker(60 * time.Second)
		// Выждем таймер для получения статистики по каждому интервалу
		<-ticker.C
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
	filter := []int{
		advertId,
	}

	bodyData, err := json.Marshal(filter)
	if err != nil {
		return promotion, fmt.Errorf("%s: ошибка маршалинга тела запроса: %w", method, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, reqURL, bytes.NewReader(bodyData))
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка создания запроса: %w", method, err)
	}

	req.Header.Set("Authorization", wb.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := wb.client.Do(req)
	if err != nil {
		return promotion, fmt.Errorf("%s: ошибка отправки запроса: %w", method, err)
	}

	if resp.StatusCode != http.StatusOK {
		return promotion, fmt.Errorf("%s: сервер ответил: %d %w", method, resp.StatusCode, entity.ErrPermanent)
	}

	var response []AdvertList
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
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

func (wb *apiClientImp) GetPromotionStats(ctx context.Context, promotions map[int64]*entity.Promotion, chunks []PromotionFilter) error {
	reqURL := fmt.Sprintf("%s%s", advertAPIURL, promoFullStatsMethod)

	bodyData, err := json.Marshal(chunks)
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
		body, _ := io.ReadAll(resp.Body)

		return fmt.Errorf("%s: сервер ответил: %d %s %w", promoFullStatsMethod, resp.StatusCode, string(body), entity.ErrPermanent)
	}

	var response []PromotionListDetail
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", promoFullStatsMethod, err)
	}

	err = resp.Body.Close()
	if err != nil {
		return fmt.Errorf("%s: ошибка закрытия тела ответа: %w", promoFullStatsMethod, err)
	}

	// Обрабатываем все полученные данные
	for _, elem := range response {
		promotion, ok := promotions[int64(elem.AdvertID)]
		if !ok {
			alogger.WarnFromCtx(ctx, "Промо %d не нашлось в карте", elem.AdvertID)
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
