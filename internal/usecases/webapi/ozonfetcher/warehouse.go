package ozonfetcher

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
)

// В Озон это пустой справочник пока пропустим его реализацию.
func (ozon *apiClientImp) GetWarehouses(ctx context.Context) ([]entity.Warehouse, error) {
	// _, _ = getWarehouses(ctx, ozon.baseURL, ozon.clientID, ozon.apiKey, ozon.metric)
	return nil, nil
}

// func getWarehouses(ctx context.Context, baseURL, clientID, apiKey string, metric metrics.Collector) ([]entity.Warehouse, error) {
// 	var timeout = time.Duration(30 * time.Second)

// 	// warehouses := []entity.Warehouse{}

// 	methodName := "/v1/warehouse/list"

// 	url := fmt.Sprintf("%s%s", baseURL, methodName)
// 	headers := make(map[string]string)
// 	headers["Client-Id"] = clientID
// 	headers["Api-Key"] = apiKey
// 	headers["Content-Type"] = "application/json"

// 	code, resp, err := httputil.SendHTTPRequest(http.MethodPost, url, []byte{}, headers, "", "", timeout)
// 	if err != nil {
// 		return nil, fmt.Errorf("%s: ошибка выполнения запроса: %w", methodName, err)
// 	}

// 	if code != http.StatusOK {
// 		return nil, fmt.Errorf("%s: ошибка выполнения запроса: %s", methodName, resp)
// 	}

// 	return nil, nil
// }
