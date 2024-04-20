package orderrepo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/cardrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/orderrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/sellerrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/warehouserepo"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

func TestOrderRepo(t *testing.T) {
	ctx := context.Background()

	conn, _, err := postgresdb.GetMockConn("../../../../migrations/data_receiver_db")
	if err != nil {
		t.Fatalf("ошибка создания мокового соединения. %s", err.Error())
	}

	// Создание продавца
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

	// Создание продавца
	sqlWarehouseRepo, err := warehouserepo.NewWarehouseRepo(ctx, conn)
	if err != nil {
		t.Fatal(err)
	}
	newWarehouse := entity.Warehouse{
		Title:    uuid.NewString(),
		ExtID:    uuid.NewString(),
		SellerID: modelSeller.ID,
	}
	modelWarehouse, err := sqlWarehouseRepo.Insert(ctx, newWarehouse)
	if err != nil {
		t.Fatal(err)
	}

	// Создание карточки
	sqlCardRepo, err := cardrepo.NewCardRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}
	newCard := entity.Card{
		VendorID:    uuid.NewString(),
		VendorCode:  uuid.NewString(),
		Title:       uuid.NewString(),
		Description: uuid.NewString(),
	}

	modelCard, err := sqlCardRepo.Insert(ctx, newCard)
	if err != nil {
		t.Fatal(err)
	}
	sqlRepo, err := orderrepo.NewOrderRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}
	newOrder := entity.Order{
		ExtID:        uuid.NewString(),
		Price:        5.5,
		Discount:     2.5,
		SpecialPrice: 10.5,
		Quantity:     5,
		Status:       uuid.NewString(),
		Type:         uuid.NewString(),
		Direction:    uuid.NewString(),
		WarehouseID:  modelWarehouse.ID,
		SellerID:     modelSeller.ID,
		CardID:       modelCard.ID,
	}
	// Создание
	model, err := sqlRepo.Insert(ctx, newOrder)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, model.ExtID, newOrder.ExtID)
	assert.Equal(t, model.Price, newOrder.Price)
  assert.Equal(t, model.Quantity, newOrder.Quantity)
	assert.Equal(t, model.Direction, newOrder.Direction)
	assert.Equal(t, model.Discount, newOrder.Discount)
	assert.Equal(t, model.SpecialPrice, newOrder.SpecialPrice)
	assert.Equal(t, model.Type, newOrder.Type)
	assert.Equal(t, model.WarehouseID, newOrder.WarehouseID)
	assert.Equal(t, model.SellerID, newOrder.SellerID)
	assert.Equal(t, model.CardID, newOrder.CardID)

	// Выборка по ID
	model, err = sqlRepo.SelectByID(ctx, model.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.ExtID, newOrder.ExtID)
	assert.Equal(t, model.Price, newOrder.Price)
  assert.Equal(t, model.Quantity, newOrder.Quantity)
	assert.Equal(t, model.Direction, newOrder.Direction)
	assert.Equal(t, model.Discount, newOrder.Discount)
	assert.Equal(t, model.SpecialPrice, newOrder.SpecialPrice)
	assert.Equal(t, model.Type, newOrder.Type)
	assert.Equal(t, model.WarehouseID, newOrder.WarehouseID)
	assert.Equal(t, model.SellerID, newOrder.SellerID)
	assert.Equal(t, model.CardID, newOrder.CardID)

	// Обновление
	newOrder.ExtID = uuid.NewString()
	newOrder.Price = 88.22
	newOrder.Discount = 55.55
	newOrder.SpecialPrice = 114.25
	newOrder.Status = uuid.NewString()
	newOrder.Type = uuid.NewString()
	newOrder.Direction = uuid.NewString()
	newOrder.ID = model.ID

	err = sqlRepo.UpdateExecOne(ctx, newOrder)
	if err != nil {
		t.Fatal(err)
	}

	// Выборка по ID
	model, err = sqlRepo.SelectByID(ctx, newOrder.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.ExtID, newOrder.ExtID)
	assert.Equal(t, model.Price, newOrder.Price)
  assert.Equal(t, model.Quantity, newOrder.Quantity)
	assert.Equal(t, model.Direction, newOrder.Direction)
	assert.Equal(t, model.Discount, newOrder.Discount)
	assert.Equal(t, model.SpecialPrice, newOrder.SpecialPrice)
	assert.Equal(t, model.Type, newOrder.Type)
	assert.Equal(t, model.WarehouseID, newOrder.WarehouseID)
	assert.Equal(t, model.SellerID, newOrder.SellerID)
	assert.Equal(t, model.CardID, newOrder.CardID)
}
