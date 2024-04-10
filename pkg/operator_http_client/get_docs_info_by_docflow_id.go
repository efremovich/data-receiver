package operator_client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (c *Client) GetDocsInfoByDocflowId(ctx context.Context, docflowId string) (*DocflowInfo, error) {
	const methodName = "OperatorClient.GetDocsInfoByDocflowId"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(`%s/system/docs_info/by_docflow_id`, c.addr), nil)
	if err != nil {
		err = fmt.Errorf("%s: ошибка создания запроса: %w", methodName, err)
		return nil, err
	}

	q := req.URL.Query()
	q.Add("docflow_id", docflowId)
	req.URL.RawQuery = q.Encode()

	if c.auth.IsRequired {
		req.SetBasicAuth(c.auth.Login, c.auth.Password)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		err = fmt.Errorf("%s: ошибка отправки запроса: %w", methodName, err)
		return nil, err
	}
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if resp.StatusCode != http.StatusOK {
		var opResp OperatorSystemResponse
		if err = json.NewDecoder(resp.Body).Decode(&opResp); err != nil {
			err = fmt.Errorf("%s: ошибка десериализации ответа в структуру ошибки: %w", methodName, err)
			return nil, err
		}

		err = fmt.Errorf("%s: ошибка получения информации о документообороте: %s", methodName, opResp.Error)
		return nil, err
	}

	var res DocflowInfo
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		err = fmt.Errorf("%s: ошибка десериализации ответа: %w", methodName, err)
		return nil, err
	}

	return &res, nil
}

type DocflowInfo struct {
	DocflowId string         `json:"docflow_id"`
	Documents []DocumentInfo `json:"documents"`
}

type DocumentInfo struct {
	Id         string    `json:"id"`
	SenderId   string    `json:"sender_id"`
	ReceiverId string    `json:"receiver_id"`
	Filename   string    `json:"filename"`
	CreatedAt  time.Time `json:"created_at"`
	Signs      []string  `json:"signs"`
}
