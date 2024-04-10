package operatorrepo

import (
	"context"
	"strings"

	"github.com/efremovich/data-receiver/internal/entity"
)

type operatorRepoMock struct {
	operators []entity.Operator
}

func NewOperatorMockRepo(operators []entity.Operator) (OperatorRepo, error) {
	return &operatorRepoMock{operators: operators}, nil
}

func (r *operatorRepoMock) GetOperators(_ context.Context) ([]entity.Operator, error) {
	return r.operators, nil
}

func (r *operatorRepoMock) GetOperatorsMap(_ context.Context) (map[string]entity.Operator, error) {
	operatorsMap := make(map[string]entity.Operator)

	for _, o := range r.operators {
		for _, thumb := range o.Thumbs {
			operatorsMap[strings.ToLower(thumb)] = o
		}
	}

	return operatorsMap, nil
}

func (r *operatorRepoMock) Ping(context.Context) error {
	return nil
}
