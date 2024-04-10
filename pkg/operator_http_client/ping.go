package operator_client

import (
	"context"
	"fmt"
	"net/http"
)

const (
	pingEndpoint = "/system/cheap_ping"
)

func (c *Client) Ping(ctx context.Context) error {
	const methodName = "OperatorClient.Ping"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/%s", c.addr, pingEndpoint), nil)
	if err != nil {
		return fmt.Errorf("%s: ошибка создания запроса: %w", methodName, err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("%s: ошибка отправки запроса: %w", methodName, err)
	}
	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%s: сервис недоступен: resp.StatusCode=%d", methodName, resp.StatusCode)
	}

	return nil
}
