package operatorfetcher

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/efremovich/data-receiver/config"
	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/pkg/logger"
	operator_client "github.com/efremovich/data-receiver/pkg/operator_http_client"
)

type OperatorFetcher interface {
	GetOperatorList(ctx context.Context) ([]entity.Operator, error)
	Ping(ctx context.Context) error
}

func New(_ context.Context, cfg config.OperatorAPI) (OperatorFetcher, error) {
	timeout := time.Second * time.Duration(cfg.TimeoutSeconds)

	c, err := operator_client.NewClient(operator_client.Config{
		Addr:               cfg.BaseURL,
		Timeout:            timeout,
		InsecureSkipVerify: false,
		Auth: operator_client.Auth{
			Login:      cfg.Login,
			Password:   cfg.Password,
			IsRequired: true,
		},
	})

	if err != nil {
		return nil, fmt.Errorf("ошибка инициализации клиента к API оператора: %w", err)
	}
	client := &opAPIclientImpl{
		client: c,
	}

	return client, nil
}

type opAPIclientImpl struct {
	client *operator_client.Client // клиент к API монолита
}

func (o *opAPIclientImpl) GetOperatorList(ctx context.Context) ([]entity.Operator, error) {
	//
	// TODO: важно, этот клиент достает только список операторов из таблицы operators
	// TODO: понять, как достать данные по Хабу, сертам ЦРПТ и пр. - доработка на операторе того же метода?
	//
	ops, err := o.client.GetOperatorsList(ctx)
	if err != nil {
		logger.GetLoggerFromContext(ctx).Errorf("ошибка получения списка операторов из клиента: %s", err.Error())
		return fallbackOperatorList, nil
	}

	res := make([]entity.Operator, 0, len(ops))

	for i := range ops {
		res = append(res, entity.Operator{
			Code:       ops[i].ID,
			Name:       ops[i].Name,
			IsDisabled: !ops[i].IsActive,
			Thumbs:     ops[i].CertificateThumbs,
		})
	}

	return res, nil
}

func (o *opAPIclientImpl) Ping(ctx context.Context) error {
	_, err := o.client.GetDocInfo(ctx, "ping_id")
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "не найден") {
			return nil
		}

		return err
	}

	return nil
}
