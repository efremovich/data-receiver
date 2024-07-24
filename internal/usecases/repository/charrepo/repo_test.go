package charrepo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/charrepo"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

func TestCharRepo(t *testing.T) {
	ctx := context.Background()

	conn, _, err := postgresdb.GetMockConn("../../../../migrations/data_base")
	if err != nil {
		t.Fatalf("ошибка создания мокового соединения. %s", err.Error())
	}

	// Создание Characteristic
	sqlCharacteristicRepo, err := charrepo.NewCharRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}

	newChar := entity.Characteristic{
		Title: uuid.NewString(),
	}

	modelChar, err := sqlCharacteristicRepo.Insert(ctx, newChar)
	if err != nil {
		t.Fatalf(err.Error())
	}

	assert.Equal(t, modelChar.Title, newChar.Title)

	modelSelectChar, err := sqlCharacteristicRepo.SelectByID(ctx, modelChar.ID)
	if err != nil {
		t.Fatalf(err.Error())
	}

	assert.Equal(t, modelChar.ID, modelSelectChar.ID)
	assert.Equal(t, modelChar.Title, modelSelectChar.Title)

	modelSelectChar, err = sqlCharacteristicRepo.SelectByTitle(ctx, modelChar.Title)
	if err != nil {
		t.Fatalf(err.Error())
	}

	assert.Equal(t, modelChar.ID, modelSelectChar.ID)
	assert.Equal(t, modelChar.Title, modelSelectChar.Title)
}

