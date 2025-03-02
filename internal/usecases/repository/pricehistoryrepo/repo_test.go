package pricehistoryrepo_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/brandrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/cardrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/pricehistoryrepo"
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
		t.Fatalf(err.Error())
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

	modelPriceSize, err := sqlPriceRepo.Insert(ctx, newPrice)
	if err != nil {
		t.Fatal(err)
	}

	sqlPriceHistoryRepo, err := pricehistoryrepo.NewPriceHistoryRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}
	newPriceHistory := entity.PriceHistory{
		Price:        2.2,
		Discount:     50000.00,
		SpecialPrice: 0,
		UpdatedAt:    time.Now(),
		PriceSizeID:  modelPriceSize.ID,
	}

	model, err := sqlPriceHistoryRepo.Insert(ctx, newPriceHistory)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, model.Price, newPriceHistory.Price)
	assert.Equal(t, model.Discount, newPriceHistory.Discount)
	assert.Equal(t, model.SpecialPrice, newPriceHistory.SpecialPrice)

	// Выборка по ID
	model, err = sqlPriceHistoryRepo.SelectByID(ctx, model.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.Price, newPriceHistory.Price)
	assert.Equal(t, model.Discount, newPriceHistory.Discount)
	assert.Equal(t, model.SpecialPrice, newPriceHistory.SpecialPrice)

	// Выборка по названию
	models, err := sqlPriceRepo.SelectByCardID(ctx, modelCard.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, model := range models {
		assert.Equal(t, model.Price, newPriceHistory.Price)
		assert.Equal(t, model.Discount, newPriceHistory.Discount)
		assert.Equal(t, model.SpecialPrice, newPriceHistory.SpecialPrice)
	}

	// Обновление
	newPriceHistory.Price = 6.6
	newPriceHistory.Discount = 1.6
	newPriceHistory.SpecialPrice = 9.0
	newPriceHistory.ID = model.ID
	err = sqlPriceHistoryRepo.UpdateExecOne(ctx, newPriceHistory)
	if err != nil {
		t.Fatal(err)
	}

	// Выборка по ID
	model, err = sqlPriceHistoryRepo.SelectByID(ctx, newPriceHistory.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.Price, newPriceHistory.Price)
	assert.Equal(t, model.Discount, newPriceHistory.Discount)
	assert.Equal(t, model.SpecialPrice, newPriceHistory.SpecialPrice)
}
