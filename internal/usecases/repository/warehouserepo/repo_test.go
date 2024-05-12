package warehouserepo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/sellerrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/warehouserepo"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

func TestWarehouseRepo(t *testing.T) {
	ctx := context.Background()

	conn, _, err := postgresdb.GetMockConn("../../../../migrations/data_receiver_db")
	if err != nil {
		t.Fatalf("ошибка создания мокового соединения. %s", err.Error())
	}

	// Создание цены
	sqlSellerRepo, err := sellerrepo.NewSellerRepo(ctx, conn)
	if err != nil {
		t.Fatal(err)
	}
	newSeller := entity.Seller{
		Title:    uuid.NewString(),
		IsEnable: true,
		ExtID:    uuid.NewString(),
	}
	modelSeller, err := sqlSellerRepo.Insert(ctx, newSeller)
	if err != nil {
		t.Fatal(err)
	}

	sqlRepo, err := warehouserepo.NewWarehouseRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}
	newWarehouse := entity.Warehouse{
		ExtID:    uuid.NewString(),
		Title:    uuid.NewString(),
		Address:  uuid.NewString(),
		Type:     uuid.NewString(),
		SellerID: modelSeller.ID,
	}
	// Создание
	model, err := sqlRepo.Insert(ctx, newWarehouse)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, model.ExtID, newWarehouse.ExtID)
	assert.Equal(t, model.Title, newWarehouse.Title)
	assert.Equal(t, model.Address, newWarehouse.Address)
	assert.Equal(t, model.Type, newWarehouse.Type)
	assert.Equal(t, model.SellerID, newWarehouse.SellerID)

	// Выборка по ID
	model, err = sqlRepo.SelectByID(ctx, model.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.ExtID, newWarehouse.ExtID)
	assert.Equal(t, model.Title, newWarehouse.Title)
	assert.Equal(t, model.Address, newWarehouse.Address)
	assert.Equal(t, model.Type, newWarehouse.Type)
	assert.Equal(t, model.SellerID, newWarehouse.SellerID)

	// Выборка по id карточки
	models, err := sqlRepo.SelectBySellerID(ctx, modelSeller.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, model := range models {
		assert.Equal(t, model.ExtID, newWarehouse.ExtID)
		assert.Equal(t, model.Title, newWarehouse.Title)
		assert.Equal(t, model.Address, newWarehouse.Address)
		assert.Equal(t, model.Type, newWarehouse.Type)
		assert.Equal(t, model.SellerID, newWarehouse.SellerID)
	}

	// Обновление
	newWarehouse.Title = uuid.NewString()
	newWarehouse.ExtID = uuid.NewString()
	newWarehouse.Address = uuid.NewString()
	newWarehouse.Type = uuid.NewString()
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

	assert.Equal(t, model.ExtID, newWarehouse.ExtID)
	assert.Equal(t, model.Title, newWarehouse.Title)
	assert.Equal(t, model.Address, newWarehouse.Address)
	assert.Equal(t, model.Type, newWarehouse.Type)
	assert.Equal(t, model.SellerID, newWarehouse.SellerID)
}
