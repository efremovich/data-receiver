package operator_client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetUserByIdEDO(ctx context.Context, idEDO string) (*User, error) {
	const methodName = "OperatorClient.GetUserByIdEDO"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(`%s/system/abonent`, c.addr), nil)
	if err != nil {
		err = fmt.Errorf("%s: ошибка создания запроса: %w", methodName, err)
		return nil, err
	}

	q := req.URL.Query()
	q.Add("id", idEDO)
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

	var res User
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		err = fmt.Errorf("%s: ошибка десериализации ответа: %w", methodName, err)
		return nil, err
	}

	return &res, nil
}

type User struct {
	ID    string `json:"guid"`             // ИдЭДО
	Inn   string `json:"inn"`              // ИНН
	Kpp   string `json:"kpp"`              // КПП
	Name  string `json:"name"`             // наименование организации
	F     string `json:"lastname"`         // Фамилия
	I     string `json:"firstname"`        // Имя
	O     string `json:"patronymic"`       // Отчество
	IsHub bool   `json:"is_hub,omitempty"` // Является ли абонентом хаба
	Type  int16  `json:"type,omitempty"`   // Тип софта
	Cert  string `json:"cert,omitempty"`   // Сертификат
}
