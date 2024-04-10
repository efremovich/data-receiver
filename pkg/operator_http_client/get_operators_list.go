package operator_client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const endpointOperators = "/system/operators_roseu?extended=true"

// OperatorInfo - модель, которую возвращает legacy API
//
//	{
//		"id": "2AD",
//		"name": "ООО Русь-Телеком",
//		"active": true,
//		"roaming": false,
//		"invs": true,
//		"direct": true,
//		"use_roseu": true,
//		"thumbs": [
//			"fd04385b9f25ff92bd7e1c5e9eee34d5aad9b3a7"
//		]
//	}
type OperatorInfo struct {
	ID                string   `json:"id"`        // код оператора ЭДО
	Name              string   `json:"name"`      // наименование оператора
	IsActive          bool     `json:"active"`    // если true - обмен с оператором включен, если false - отключен
	IsRoaming         bool     `json:"roaming"`   // если true - это роуминговый оператор
	UsesInvitations   bool     `json:"invs"`      // если true - поддерживает приглашения по Технологии ФНС
	IsDirect          bool     `json:"direct"`    // если true - поддерживает прямой обмен (без участия 1С-Хаб)
	UsesFNSTech       bool     `json:"use_roseu"` // поддерживает ли обмен по Технологии ФНС
	CertificateThumbs []string `json:"thumbs"`
}

func (c *Client) GetOperatorsList(ctx context.Context) ([]OperatorInfo, error) {
	const methodName = "OperatorClient.GetOperatorsList"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s%s", c.addr, endpointOperators), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка создания запроса: %w", methodName, err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка отправки запроса: %w", methodName, err)
	}
	defer resp.Body.Close()

	var opMap map[string]OperatorInfo

	if err := json.NewDecoder(resp.Body).Decode(&opMap); err != nil {
		return nil, fmt.Errorf("%s: ошибка чтения/десериализации тела ответа: %w", methodName, err)
	}

	var operatorsList []OperatorInfo

	for _, v := range opMap {
		operatorsList = append(operatorsList, v)
	}

	return operatorsList, nil
}
