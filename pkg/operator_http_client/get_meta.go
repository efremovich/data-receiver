package operator_client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetMetaJson(ctx context.Context, docId string) (*MetaJson, error) {
	const methodName = "OperatorClient.GetMetaJson"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(`%s/system/file`, c.addr), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка создания запроса: %w", methodName, err)
	}

	q := req.URL.Query()
	q.Add("parent_id", docId)
	q.Add("type", fileTypeMeta)
	req.URL.RawQuery = q.Encode()

	if c.auth.IsRequired {
		req.SetBasicAuth(c.auth.Login, c.auth.Password)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка отправки запроса: %w", methodName, err)
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

		err = fmt.Errorf("%s: ошибка получения информации о документе: %s", methodName, opResp.Error)
		return nil, err
	}

	var res MetaJson
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		err = fmt.Errorf("%s: ошибка десериализации ответа: %w", methodName, err)
		return nil, err
	}

	return &res, nil
}

type MetaJson struct {
	ID       string `yaml:"id" json:"ID"`
	Sender   string `yaml:"sender" json:"Sender,omitempty"`
	Receiver string `yaml:"receiver" json:"Receiver,omitempty"`
	History  string `yaml:"history" json:"History,omitempty"`
	Signs    []*struct {
		ID string `yaml:"id"`
	} `yaml:"signs" json:"Signs,omitempty"`
}
