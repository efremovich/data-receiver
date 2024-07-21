package mediafiletypeenumrepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
)

type mediaFileTypeEnumDB struct {
	ID   int64  `db:"id"`
	Type string `db:"type"`
}

func convertToDBMediaFileTypeEnum(_ context.Context, in entity.MediaFileTypeEnum) *mediaFileTypeEnumDB {
	return &mediaFileTypeEnumDB{
		ID:   in.ID,
		Type: in.Type,
	}
}

func (c mediaFileTypeEnumDB) ConvertToEntityMediaFileTypeEnum(_ context.Context) *entity.MediaFileTypeEnum {
	return &entity.MediaFileTypeEnum{
		ID:   c.ID,
		Type: c.Type,
	}
}
