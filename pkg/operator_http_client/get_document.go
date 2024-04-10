package operator_client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetDocument(ctx context.Context, docId string) (*bytes.Buffer, error) {
	const methodName = "OperatorClient.GetDocument"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(`%s/system/file`, c.addr), nil)
	if err != nil {
		return nil, fmt.Errorf("%s: ошибка создания запроса: %w", methodName, err)
	}

	q := req.URL.Query()
	q.Add("id", docId)
	q.Add("type", fileTypeDocument)
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

	var buff bytes.Buffer
	if _, err = buff.ReadFrom(resp.Body); err != nil {
		return nil, fmt.Errorf("%s: ошибка чтения тела ответа: %w", methodName, err)
	}

	if resp.StatusCode != http.StatusOK {
		var opResp OperatorSystemResponse
		if err = json.Unmarshal(buff.Bytes(), &opResp); err != nil {
			return nil, fmt.Errorf("%s: ошибка десериализации ответа в структуру ошибки: %w", methodName, err)
		}

		return nil, fmt.Errorf("%s: ошибка получения файла подписи: %s", methodName, opResp.Error)
	}

	return &buff, nil
}
