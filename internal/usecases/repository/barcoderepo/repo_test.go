package barcoderepo_test

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
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

func TestBarcodeRepo(t *testing.T) {
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

	// Создание бренда
	sqlBrandRepo, err := brandrepo.NewBrandRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}
	newBrand := entity.Brand{
		ExternalID: 1,
		Title:      uuid.NewString(),
		SellerID:   modelSeller.ID,
	}
	modelBrand, err := sqlBrandRepo.Insert(ctx, newBrand)
	if err != nil {
		t.Fatal(err)
	}

	// Создание карточки
	sqlCardRepo, err := cardrepo.NewCardRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}
	newCard := entity.Card{
		Brand:       *modelBrand,
		VendorID:    uuid.NewString(),
		VendorCode:  uuid.NewString(),
		Title:       uuid.NewString(),
		Description: uuid.NewString(),
	}

	modelCard, err := sqlCardRepo.Insert(ctx, newCard)
	if err != nil {
		t.Fatal(err)
	}

	// Создание размера
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
		SizeID:       modelSize.ID,
		CardID:       modelCard.ID,
	}

	modelPrice, err := sqlPriceRepo.Insert(ctx, newPrice)
	if err != nil {
		t.Fatal(err)
	}

	sqlRepo, err := barcoderepo.NewBarcodeRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}
	newBarcode := entity.Barcode{
		Barcode:     uuid.NewString(),
		PriceSizeID: modelPrice.ID,
		SellerID:    modelSeller.ID,
	}
	// Создание
	model, err := sqlRepo.Insert(ctx, newBarcode)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.Barcode, newBarcode.Barcode)

	// Выборка по ID
	model, err = sqlRepo.SelectByBarcode(ctx, model.Barcode)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.Barcode, newBarcode.Barcode)

	// Обновление
	newBarcode.Barcode = model.Barcode

	err = sqlRepo.UpdateExecOne(ctx, newBarcode)
	if err != nil {
		t.Fatal(err)
	}

	// Выборка по ID
	model, err = sqlRepo.SelectByBarcode(ctx, newBarcode.Barcode)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, model.Barcode, newBarcode.Barcode)
}
