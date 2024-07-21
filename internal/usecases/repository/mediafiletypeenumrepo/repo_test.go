package mediafiletypeenumrepo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/mediafiletypeenumrepo"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

func TestMediaFileRepo(t *testing.T) {
	ctx := context.Background()

	conn, _, err := postgresdb.GetMockConn("../../../../migrations/data_base")
	if err != nil {
		t.Fatalf("ошибка создания мокового соединения. %s", err.Error())
	}

	sqlMediaFileTypeEnumRepo, err := mediafiletypeenumrepo.NewMediaFileTypeEnumRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}

	newMediaFileTypeEnum := entity.MediaFileTypeEnum{
		Type: "VIDEO",
	}

	// Создание
	modelMediaFileTypeEnum, err := sqlMediaFileTypeEnumRepo.Insert(ctx, newMediaFileTypeEnum)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, newMediaFileTypeEnum.Type, modelMediaFileTypeEnum.Type)

	// Выборка по ID
	modelSelectMediaFileTypeEnum, err := sqlMediaFileTypeEnumRepo.SelectByID(ctx, modelMediaFileTypeEnum.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, newMediaFileTypeEnum.Type, modelSelectMediaFileTypeEnum.Type)

	// Обновление
	newMediaFileTypeEnum.ID = modelSelectMediaFileTypeEnum.ID
	newMediaFileTypeEnum.Type = uuid.NewString()

	err = sqlMediaFileTypeEnumRepo.UpdateExecOne(ctx, newMediaFileTypeEnum)
	if err != nil {
		t.Fatal(err)
	}

	// Выборка по Type
  model, err := sqlMediaFileTypeEnumRepo.SelectByType(ctx, modelSelectMediaFileTypeEnum.Type)
	if err != nil {
		t.Fatal(err)
	}
		assert.Equal(t, modelMediaFileTypeEnum.Type, model.Type)
}
