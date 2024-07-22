package countryrepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
)

type countryDB struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

func convertToDBCountry(_ context.Context, in entity.Country) *countryDB {
	return &countryDB{
		ID:   in.ID,
		Name: in.Name,
	}
}

func (c countryDB) ConvertToEntityCountry(_ context.Context) *entity.Country {
	return &entity.Country{
		ID:   c.ID,
		Name: c.Name,
	}
}
