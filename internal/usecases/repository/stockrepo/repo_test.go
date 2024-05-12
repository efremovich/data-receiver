package stockrepo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/barcoderepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/cardrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/pricerepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/sellerrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/sizerepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/stockrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/warehouserepo"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

func TestStockRepo(t *testing.T) {
	ctx := context.Background()

	conn, _, err := postgresdb.GetMockConn("../../../../migrations/data_receiver_db")
	if err != nil {
		t.Fatalf("ошибка создания мокового соединения. %s", err.Error())
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

	// Создание цены
	sqlPriceRepo, err := pricerepo.NewPriceRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}
	newPrice := entity.Price{
		Price:        5.5,
		Discount:     1.5,
		SpecialPrice: 8.0,
		SellerID:     modelSeller.ID,
		CardID:       modelCard.ID,
	}

	modelPrice, err := sqlPriceRepo.Insert(ctx, newPrice)
	if err != nil {
		t.Fatal(err)
	}

	sqlSizeRepo, err := sizerepo.NewSizeRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}
	newSize := entity.Size{
		TechSize: uuid.NewString(),
		Title:    uuid.NewString(),
		CardID:   modelCard.ID,
		PriceID:  modelPrice.ID,
	}
	// Создание
	modelSize, err := sqlSizeRepo.Insert(ctx, newSize)
	if err != nil {
		t.Fatal(err)
	}

	sqlBarcodeRepo, err := barcoderepo.NewBarcodeRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}
	newBarcode := entity.Barcode{
		Barcode:  uuid.NewString(),
		SizeID:   modelSize.ID,
		SellerID: modelSeller.ID,
	}
	modelBarcode, err := sqlBarcodeRepo.Insert(ctx, newBarcode)
	if err != nil {
		t.Fatal(err)
	}

	sqlWarehouseRepo, err := warehouserepo.NewWarehouseRepo(ctx, conn)
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
	modelWarehouse, err := sqlWarehouseRepo.Insert(ctx, newWarehouse)
	if err != nil {
		t.Fatal(err)
	}

	sqlRepo, err := stockrepo.NewStockRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}
	newStock := entity.Stock{
		Quantity:        55,
		InWayToClient:   500,
		InWayFromClient: 21,
		SizeID:          modelSize.ID,
		Barcode:         modelBarcode.Barcode,
		WarehouseID:     modelWarehouse.ID,
		CardID:          modelCard.ID,
		SellerID:        modelSeller.ID,
	}
	// Создание
	model, err := sqlRepo.Insert(ctx, newStock)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, model.Quantity, newStock.Quantity)
	assert.Equal(t, model.InWayFromClient, newStock.InWayFromClient)
	assert.Equal(t, model.InWayToClient, newStock.InWayToClient)
	assert.Equal(t, model.SizeID, newStock.SizeID)
	assert.Equal(t, model.Barcode, newStock.Barcode)
	assert.Equal(t, model.WarehouseID, newStock.WarehouseID)
	assert.Equal(t, model.CardID, newStock.CardID)
	assert.Equal(t, model.SellerID, newStock.SellerID)

	// Выборка по ID
	model, err = sqlRepo.SelectByID(ctx, model.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.Quantity, newStock.Quantity)
	assert.Equal(t, model.InWayFromClient, newStock.InWayFromClient)
	assert.Equal(t, model.InWayToClient, newStock.InWayToClient)
	assert.Equal(t, model.SizeID, newStock.SizeID)
	assert.Equal(t, model.Barcode, newStock.Barcode)
	assert.Equal(t, model.WarehouseID, newStock.WarehouseID)
	assert.Equal(t, model.CardID, newStock.CardID)
	assert.Equal(t, model.SellerID, newStock.SellerID)
	assert.Equal(t, model.SellerID, newStock.SellerID)

	// Выборка по id карточки
	models, err := sqlRepo.SelectBySellerID(ctx, modelSeller.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, model := range models {
		assert.Equal(t, model.Quantity, newStock.Quantity)
		assert.Equal(t, model.InWayFromClient, newStock.InWayFromClient)
		assert.Equal(t, model.InWayToClient, newStock.InWayToClient)
		assert.Equal(t, model.SizeID, newStock.SizeID)
		assert.Equal(t, model.Barcode, newStock.Barcode)
		assert.Equal(t, model.WarehouseID, newStock.WarehouseID)
		assert.Equal(t, model.CardID, newStock.CardID)
		assert.Equal(t, model.SellerID, newStock.SellerID)
	}

	// Обновление
	newStock.Quantity = 10
	newStock.InWayToClient = 546
	newStock.InWayFromClient = 0
	newStock.SellerID = modelSeller.ID
	newStock.ID = model.ID

	err = sqlRepo.UpdateExecOne(ctx, newStock)
	if err != nil {
		t.Fatal(err)
	}

	// Выборка по ID
	model, err = sqlRepo.SelectByID(ctx, newStock.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.Quantity, newStock.Quantity)
	assert.Equal(t, model.InWayFromClient, newStock.InWayFromClient)
	assert.Equal(t, model.InWayToClient, newStock.InWayToClient)
	assert.Equal(t, model.SizeID, newStock.SizeID)
	assert.Equal(t, model.Barcode, newStock.Barcode)
	assert.Equal(t, model.WarehouseID, newStock.WarehouseID)
	assert.Equal(t, model.CardID, newStock.CardID)
	assert.Equal(t, model.SellerID, newStock.SellerID)
}
