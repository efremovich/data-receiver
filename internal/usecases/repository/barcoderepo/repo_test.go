package barcoderepo_test

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
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

func TestBarcodeRepo(t *testing.T) {
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

	sqlRepo, err := barcoderepo.NewBarcodeRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}
	newBarcode := entity.Barcode{
		Barcode:  uuid.NewString(),
		SizeID:   modelSize.ID,
		SellerID: modelSeller.ID,
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