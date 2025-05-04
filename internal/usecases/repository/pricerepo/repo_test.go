package pricerepo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/brandrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/cardrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/pricerepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/sellerrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/sizerepo"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

func TestPriceRepo(t *testing.T) {
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
	sqlBrandRepo, err := brandrepo.NewBrandRepo(ctx, conn)
	if err != nil {
		t.Fatal(err.Error())
	}
	newBrand := entity.Brand{
		Title:    uuid.NewString(),
		SellerID: modelSeller.ID,
	}
	modelBrand, err := sqlBrandRepo.Insert(ctx, newBrand)
	if err != nil {
		t.Fatal(err)
	}
	// Создание карточки
	sqlCardRepo, err := cardrepo.NewCardRepo(ctx, conn)
	if err != nil {
		t.Fatal(err.Error())
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

	sqlSizeRepo, err := sizerepo.NewSizeRepo(ctx, conn)
	if err != nil {
		t.Fatal(err.Error())
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
		t.Fatal(err.Error())
	}
	newPrice := entity.PriceSize{
		Price:                5.5,
		PriceWithoutDiscount: 1.5,
		PriceFinal:           8.0,
		CardID:               modelCard.ID,
		SizeID:               modelSize.ID,
	}

	// Создание
	model, err := sqlPriceRepo.Insert(ctx, &newPrice)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, model.Price, newPrice.Price)
	assert.Equal(t, model.PriceWithoutDiscount, newPrice.PriceWithoutDiscount)
	assert.Equal(t, model.PriceFinal, newPrice.PriceFinal)

	// Выборка по ID
	model, err = sqlPriceRepo.SelectByID(ctx, model.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.Price, newPrice.Price)
	assert.Equal(t, model.PriceWithoutDiscount, newPrice.PriceWithoutDiscount)
	assert.Equal(t, model.PriceFinal, newPrice.PriceFinal)

	// Выборка по названию
	models, err := sqlPriceRepo.SelectByCardID(ctx, modelCard.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, model := range models {
		assert.Equal(t, model.Price, newPrice.Price)
		assert.Equal(t, model.PriceWithoutDiscount, newPrice.PriceWithoutDiscount)
		assert.Equal(t, model.PriceFinal, newPrice.PriceFinal)
	}

	// Обновление
	newPrice.Price = 6.6
	newPrice.PriceWithoutDiscount = 1.6
	newPrice.PriceFinal = 9.0
	newPrice.ID = model.ID
	newPrice.CardID = modelCard.ID
	newPrice.SizeID = modelSize.ID

	err = sqlPriceRepo.UpdateExecOne(ctx, &newPrice)
	if err != nil {
		t.Fatal(err)
	}

	// Выборка по ID
	model, err = sqlPriceRepo.SelectByID(ctx, newPrice.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.Price, newPrice.Price)
	assert.Equal(t, model.PriceWithoutDiscount, newPrice.PriceWithoutDiscount)
	assert.Equal(t, model.PriceFinal, newPrice.PriceFinal)
}
