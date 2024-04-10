package operator_client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (c *Client) GetDocsLinksByDocId(ctx context.Context, docId string) (*DocLinks, error) {
	const methodName = "OperatorClient.GetDocsLinksByDocId"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(`%s/system/doc_links/by_doc_id`, c.addr), nil)
	if err != nil {
		err = fmt.Errorf("%s: ошибка создания запроса: %w", methodName, err)
		return nil, err
	}

	q := req.URL.Query()
	q.Add("doc_id", docId)
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

		err = fmt.Errorf("%s: ошибка получения информации о документе: %s", methodName, opResp.Error)
		return nil, err
	}

	var res DocLinks
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		err = fmt.Errorf("%s: ошибка десериализации ответа: %w", methodName, err)
		return nil, err
	}

	return &res, nil
}

type DocLinks struct {
	DocId     string         `json:"doc_id"`
	DocflowId string         `json:"docflow_id"`
	LSs       []LS           `json:"ls"`
	EDIs      []EDIContainer `json:"edi"`
}

// GetFirstInboxTP - Вернет первый входящий ТП, или nil
func (d *DocLinks) GetFirstInboxTP() *TP {
	if len(d.LSs) != 0 && d.LSs[1].TP.Id != "" {
		tp := d.LSs[1].TP
		return &tp
	}

	return nil
}

type LS struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	IsInbox   bool      `json:"is_inbox"`
	TP        TP        `json:"tp"`
}

type TP struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	IsInbox   bool      `json:"is_inbox"`
}

type EDIContainer struct {
	Id        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	IsInbox   bool      `json:"is_inbox"`
}
