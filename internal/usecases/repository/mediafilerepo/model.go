package mediafilerepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
)

type mediaFileDB struct {
	ID   int64  `db:"id"`
	Link string `db:"link"`

	CardID int64 `db:"card_id"`
}

func convertToDBMediaFile(_ context.Context, in entity.MediaFile) *mediaFileDB {
	return &mediaFileDB{
		ID:   in.ID,
		Link: in.Link,

		CardID: in.CardID,
	}
}

func (c mediaFileDB) ConvertToEntityMediaFile(_ context.Context) *entity.MediaFile {
	return &entity.MediaFile{
		ID:     c.ID,
		Link:   c.Link,
		CardID: c.CardID,
	}
}
