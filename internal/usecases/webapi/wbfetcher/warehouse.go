package wbfetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/efremovich/data-receiver/internal/entity"
)

var warehouseTypes = map[int]string{
	1: "обычный",
	2: "СГТ (Сверхгабаритный товар)",
	3: "КГТ",
}

type WarehouseResponce struct {
	Address      string  `json:"address"`
	Name         string  `json:"name"`
	City         string  `json:"city"`
	ID           int     `json:"id"`
	Longitude    float64 `json:"longitude"`
	Latitude     float64 `json:"latitude"`
	CargoType    int     `json:"cargoType"`
	DeliveryType int     `json:"deliveryType"`
	Selected     bool    `json:"selected"`
}

func (wb *apiClientImp) GetWarehouses(ctx context.Context) ([]entity.Warehouse, error) {
	const methodName = "/api/v3/offices"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s", marketPlaceAPIURL, methodName), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка создания запроса: %w", methodName, err)
	}

	req.Header.Set("Authorization", wb.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept", "application/json")

	resp, err := wb.client.Do(req)

	if err != nil && resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: ошибка отправки запроса: %w", methodName, err)
	}
	defer resp.Body.Close()

	var responce []WarehouseResponce
	if err := json.NewDecoder(resp.Body).Decode(&responce); err != nil {
		return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", methodName, err)
	}

	var warehouses []entity.Warehouse

	for _, elem := range responce {
		warehouse := entity.Warehouse{
			ExternalID: int64(elem.ID),
			Title:      elem.Name,
			Address:    elem.Address,
			TypeName:   warehouseTypes[elem.CargoType],
		}

		warehouses = append(warehouses, warehouse)
	}

	return warehouses, nil
}
