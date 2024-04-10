package operator_client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// GetDocInfo - загрузка информации о документе из БД оператора (почти все поля из таблицы docs)
func (c *Client) GetDocInfo(ctx context.Context, docId string) (*DocInfo, error) {
	const methodName = "OperatorClient.GetDocInfo"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(`%s/system/doc_info/by_id`, c.addr), nil)
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

	var res DocInfo
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		err = fmt.Errorf("%s: ошибка десериализации ответа: %w", methodName, err)
		return nil, err
	}

	return &res, nil
}

type DocInfo struct {
	Id           string    `json:"id"`            // Id документа
	Route        string    `json:"route"`         // Маршрут документа
	History      string    `json:"history"`       // История и тип документа
	DocflowId    string    `json:"docflow_id"`    // Id документооборота
	Sender       string    `json:"sender_id"`     // ИдЭДО отправителя документа
	Receiver     string    `json:"receiver_id"`   // ИдЭДО получателя документа
	Filename     string    `json:"filename"`      // Имя файла документа с расширением
	SignedAt     time.Time `json:"signed_at"`     // Дата формирования документа отправителем
	CreatedAt    time.Time `json:"created_at"`    // Дата получения документа оператором
	ResolvedAt   time.Time `json:"resolved_at"`   // Момент завершения проверок
	ProcessedAt  time.Time `json:"processed_at"`  // Момент готовности обработки
	ConfirmedAt  time.Time `json:"confirmed_at"`  // Момент формирования подтверждений
	IsNeedSign   bool      `json:"is_need_sign"`  // Флаг необходимости ответной подписи
	IsEncrypted  bool      `json:"is_encrypted"`  // Зашифровано содержимое документа или нет
	IsCompressed bool      `json:"is_compressed"` // Сжато содержимое документа или нет
	ParentsIds   []string  `json:"parents_ids"`   // Список Id родительских документов. Пустой, если документ первичный
	ErrCode      string    `json:"err_code"`      // Код ошибки
	ErrText      string    `json:"err_text"`      // Текст ошибки
}
