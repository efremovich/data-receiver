package operator_client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetDocsInfoByDocflowIdFast(ctx context.Context, docflowId string) (*DocflowInfoFast, error) {
	const methodName = "OperatorClient.GetDocsInfoByDocflowIdFast"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(`%s/system/docs_info/by_docflow_id`, c.addr), nil)
	if err != nil {
		err = fmt.Errorf("%s: ошибка создания запроса: %w", methodName, err)
		return nil, err
	}

	q := req.URL.Query()
	q.Add("docflow_id", docflowId)
	q.Add("fast", "true")
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

	var res DocflowInfoFast
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		err = fmt.Errorf("%s: ошибка десериализации ответа: %w", methodName, err)
		return nil, err
	}

	return &res, nil
}

type DocflowInfoFast struct {
	DocflowId string              `json:"docflow_id"`
	Documents []*DocumentInfoFast `json:"documents"`
}

type DocumentInfoFast struct {
	Id      string `json:"id"`
	History string `json:"history"`
}
