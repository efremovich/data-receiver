package statusrepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
)

type statusDB struct {
	ID   int64
	Name string
}

func convertToDBStatus(_ context.Context, in entity.Status) *statusDB {
	return &statusDB{
		ID:   in.ID,
		Name: in.Name,
	}
}

func (s *statusDB) convertToEntityStatus(_ context.Context) *entity.Status {
	return &entity.Status{
		ID:   s.ID,
		Name: s.Name,
	}
}
