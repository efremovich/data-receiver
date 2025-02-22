package stockrepo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/barcoderepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/brandrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/cardrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/pricerepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/sellerrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/sizerepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/stockrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/warehouserepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/warehousetyperepo"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

func TestStockRepo(t *testing.T) {
	ctx := context.Background()

	conn, _, err := postgresdb.GetMockConn("../../../../migrations/data_base")
	if err != nil {
		t.Fatalf("ошибка создания мокового соединения. %s", err.Error())
	}
	// Создание Seller
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

	// Создание Brand
	sqlRepo, err := brandrepo.NewBrandRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}

	newBrand := entity.Brand{
		Title:    uuid.NewString(),
		SellerID: modelSeller.ID,
	}
	modelBrand, err := sqlRepo.Insert(ctx, newBrand)
	if err != nil {
		t.Fatal(err)
	}
	// Создание карточки
	sqlCardRepo, err := cardrepo.NewCardRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}
	newCard := entity.Card{
		ExternalID:  0,
		VendorID:    uuid.NewString(),
		VendorCode:  uuid.NewString(),
		Title:       uuid.NewString(),
		Description: uuid.NewString(),
		Brand:       *modelBrand,
	}

	modelCard, err := sqlCardRepo.Insert(ctx, newCard)
	if err != nil {
		t.Fatal(err)
	}
	// Создание Size
	sqlSizeRepo, err := sizerepo.NewSizeRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}
	newSize := entity.Size{
		TechSize: uuid.NewString(),
		Title:    uuid.NewString(),
	}
	// Создание
	modelSize, err := sqlSizeRepo.Insert(ctx, newSize)
	if err != nil {
		t.Fatal(err)
	}

	// Создание цены
	sqlPriceRepo, err := pricerepo.NewPriceRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}
	newPrice := entity.PriceSize{
		Price:        5.5,
		Discount:     1.5,
		SpecialPrice: 8.0,
		CardID:       modelCard.ID,
		SizeID:       modelSize.ID,
	}

	modelPrice, err := sqlPriceRepo.Insert(ctx, newPrice)
	if err != nil {
		t.Fatal(err)
	}

	sqlBarcodeRepo, err := barcoderepo.NewBarcodeRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}
	newBarcode := entity.Barcode{
		Barcode:     uuid.NewString(),
		PriceSizeID: modelPrice.ID,
		SellerID:    modelSeller.ID,
	}
	modelBarcode, err := sqlBarcodeRepo.Insert(ctx, newBarcode)
	if err != nil {
		t.Fatal(err)
	}

	sqlWarehouseTypeRepo, err := warehousetyperepo.NewWarehouseTypeRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}
	newWarehouseType := entity.WarehouseType{
		Title: uuid.NewString(),
	}
	modelWarehouseType, err := sqlWarehouseTypeRepo.Insert(ctx, newWarehouseType)
	if err != nil {
		t.Fatal(err)
	}

	// Создание warehouse
	sqlWarehouseRepo, err := warehouserepo.NewWarehouseRepo(ctx, conn)
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
	modelWarehouse, err := sqlWarehouseRepo.Insert(ctx, newWarehouse)
	if err != nil {
		t.Fatal(err)
	}

	sqlStockRepo, err := stockrepo.NewStockRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}
	newStock := entity.Stock{
		Quantity:    55,
		BarcodeID:   modelBarcode.ID,
		WarehouseID: modelWarehouse.ID,
		CardID:      modelCard.ID,
	}
	// Создание
	model, err := sqlStockRepo.Insert(ctx, newStock)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, model.Quantity, newStock.Quantity)
	assert.Equal(t, model.BarcodeID, newStock.BarcodeID)
	assert.Equal(t, model.WarehouseID, newStock.WarehouseID)
	assert.Equal(t, model.CardID, newStock.CardID)

	// Выборка по ID
	model, err = sqlStockRepo.SelectByID(ctx, model.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.Quantity, newStock.Quantity)
	assert.Equal(t, model.BarcodeID, newStock.BarcodeID)
	assert.Equal(t, model.WarehouseID, newStock.WarehouseID)
	assert.Equal(t, model.CardID, newStock.CardID)

	// Обновление
	newStock.Quantity = 10
	newStock.SellerID = modelSeller.ID
	newStock.ID = model.ID

	err = sqlStockRepo.UpdateExecOne(ctx, newStock)
	if err != nil {
		t.Fatal(err)
	}
}
