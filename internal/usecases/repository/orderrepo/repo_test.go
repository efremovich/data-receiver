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

	conn, _, err := postgresdb.GetMockConn("../../../../migrations/data_base")
	if err != nil {
		t.Fatalf("ошибка создания мокового соединения. %s", err.Error())
	}

	// Создание продавца
	sqlSellerRepo, err := sellerrepo.NewSellerRepo(ctx, conn)
	if err != nil {
		t.Fatal(err)
	}
	newSeller := entity.MarketPlace{
		Title:      uuid.NewString(),
		IsEnabled:  true,
		ExternalID: uuid.NewString(),
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
		Title:      uuid.NewString(),
		ExternalID: 0,
		SellerID:   modelSeller.ID,
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
		ID:         0,
		ExternalID: uuid.NewString(),
		Price:      5.5,
		Type:       uuid.NewString(),
		Direction:  uuid.NewString(),
		Sale:       0,
		Quantity:   5,
		Status:     &entity.Status{},
		Region:     &entity.Region{},
		Warehouse:  modelWarehouse,
		Seller:     modelSeller,
		Card:       modelCard,
	}
	// Создание
	model, err := sqlRepo.Insert(ctx, newOrder)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, model.ExternalID, newOrder.ExternalID)
	assert.Equal(t, model.Price, newOrder.Price)
	assert.Equal(t, model.Quantity, newOrder.Quantity)
	assert.Equal(t, model.Direction, newOrder.Direction)

	// Выборка по ID
	model, err = sqlRepo.SelectByID(ctx, model.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.ExternalID, newOrder.ExternalID)
	assert.Equal(t, model.Price, newOrder.Price)
	assert.Equal(t, model.Quantity, newOrder.Quantity)
	assert.Equal(t, model.Direction, newOrder.Direction)

	// Обновление
	newOrder.ExternalID = uuid.NewString()
	newOrder.Price = 88.22

	err = sqlRepo.UpdateExecOne(ctx, newOrder)
	if err != nil {
		t.Fatal(err)
	}

	// Выборка по ID
	model, err = sqlRepo.SelectByID(ctx, newOrder.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.ExternalID, newOrder.ExternalID)
	assert.Equal(t, model.Price, newOrder.Price)
	assert.Equal(t, model.Quantity, newOrder.Quantity)
	assert.Equal(t, model.Direction, newOrder.Direction)
}
