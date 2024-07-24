package warehousetyperepo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/warehousetyperepo"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

func TestWarehousetypeRepo(t *testing.T) {
	ctx := context.Background()

	conn, _, err := postgresdb.GetMockConn("../../../../migrations/data_base")
	if err != nil {
		t.Fatalf("ошибка создания мокового соединения. %s", err.Error())
	}

	// Создание WarehouseType
	sqlWarehouseRepo, err := warehousetyperepo.NewWarehouseTypeRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}
	newWarehouseType := entity.WarehouseType{
		Title: uuid.NewString(),
	}
	// Создание
	modelWarehouseType, err := sqlWarehouseRepo.Insert(ctx, newWarehouseType)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, modelWarehouseType.Title, newWarehouseType.Title)

	modelSelectByIDWarehouseType, err := sqlWarehouseRepo.SelectByID(ctx, modelWarehouseType.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, modelWarehouseType.ID, modelSelectByIDWarehouseType.ID)

	modelSelectByTitleWarehouseType, err := sqlWarehouseRepo.SelectByTitle(ctx, modelWarehouseType.Title)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, modelWarehouseType.Title, modelSelectByTitleWarehouseType.Title)

	err = sqlWarehouseRepo.UpdateExecOne(ctx, *modelSelectByIDWarehouseType)
	if err != nil {
		t.Fatal(err)
	}
}
