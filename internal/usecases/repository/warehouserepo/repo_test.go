package warehouserepo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/sellerrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/warehouserepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/warehousetyperepo"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

func TestWarehouseRepo(t *testing.T) {
	ctx := context.Background()

	conn, _, err := postgresdb.GetMockConn("../../../../migrations/data_base")
	if err != nil {
		t.Fatalf("ошибка создания мокового соединения. %s", err.Error())
	}

	// Создание цены
	sqlSellerRepo, err := sellerrepo.NewSellerRepo(ctx, conn)
	if err != nil {
		t.Fatal(err)
	}
	newSeller := entity.Seller{
		Title:      uuid.NewString(),
		IsEnabled:  true,
		ExternalID: uuid.NewString(),
	}
	modelSeller, err := sqlSellerRepo.Insert(ctx, newSeller)
	if err != nil {
		t.Fatal(err)
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
	sqlRepo, err := warehouserepo.NewWarehouseRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}
	newWarehouse := entity.Warehouse{
		ExternalID: 0,
		Title:      uuid.NewString(),
		Address:    uuid.NewString(),
		TypeID:     modelWarehouseType.ID,
		SellerID:   modelSeller.ID,
	}
	// Создание
	model, err := sqlRepo.Insert(ctx, newWarehouse)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, model.ExternalID, newWarehouse.ExternalID)
	assert.Equal(t, model.Title, newWarehouse.Title)
	assert.Equal(t, model.Address, newWarehouse.Address)
	assert.Equal(t, model.TypeID, newWarehouse.TypeID)
	assert.Equal(t, model.SellerID, newWarehouse.SellerID)

	// Выборка по ID
	model, err = sqlRepo.SelectByID(ctx, model.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.ExternalID, newWarehouse.ExternalID)
	assert.Equal(t, model.Title, newWarehouse.Title)
	assert.Equal(t, model.Address, newWarehouse.Address)
	assert.Equal(t, model.TypeID, newWarehouse.TypeID)
	assert.Equal(t, model.SellerID, newWarehouse.SellerID)

	// Выборка по id карточки
	models, err := sqlRepo.SelectBySellerID(ctx, modelSeller.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, model := range models {
		assert.Equal(t, model.ExternalID, newWarehouse.ExternalID)
		assert.Equal(t, model.Title, newWarehouse.Title)
		assert.Equal(t, model.Address, newWarehouse.Address)
		assert.Equal(t, model.TypeID, newWarehouse.TypeID)
		assert.Equal(t, model.SellerID, newWarehouse.SellerID)
	}

	// Обновление
	newWarehouse.Title = uuid.NewString()
	newWarehouse.ExternalID = 2
	newWarehouse.Address = uuid.NewString()
	newWarehouse.TypeID = model.TypeID
	newWarehouse.SellerID = modelSeller.ID
	newWarehouse.ID = model.ID

	err = sqlRepo.UpdateExecOne(ctx, newWarehouse)
	if err != nil {
		t.Fatal(err)
	}

	// Выборка по ID
	model, err = sqlRepo.SelectByID(ctx, newWarehouse.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.ExternalID, newWarehouse.ExternalID)
	assert.Equal(t, model.Title, newWarehouse.Title)
	assert.Equal(t, model.Address, newWarehouse.Address)
	assert.Equal(t, model.TypeID, newWarehouse.TypeID)
	assert.Equal(t, model.SellerID, newWarehouse.SellerID)
}
