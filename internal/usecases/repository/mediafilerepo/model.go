package mediafilerepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
)

type mediaFileDB struct {
	ID   int64  `db:"id"`
	Link string `db:"link"`

  TypeID int64 `db:"type_id"`
	CardID int64 `db:"card_id"`
}

func convertToDBMediaFile(_ context.Context, in entity.MediaFile) *mediaFileDB {
	return &mediaFileDB{
		ID:     in.ID,
		Link:   in.Link,
		TypeID: in.TypeID,
		CardID: in.CardID,
	}
}

func (c mediaFileDB) ConvertToEntityMediaFile(_ context.Context) *entity.MediaFile {
	return &entity.MediaFile{
		ID:     c.ID,
		Link:   c.Link,
		CardID: c.CardID,
		TypeID: c.TypeID,
	}
}

