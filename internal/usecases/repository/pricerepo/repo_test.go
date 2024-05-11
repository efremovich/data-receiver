package pricerepo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/cardrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/pricerepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/sellerrepo"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

func TestPriceRepo(t *testing.T) {
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

	sqlRepo, err := pricerepo.NewPriceRepo(ctx, conn)
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
	// Создание
	model, err := sqlRepo.Insert(ctx, newPrice)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, model.Price, newPrice.Price)
	assert.Equal(t, model.Discount, newPrice.Discount)
	assert.Equal(t, model.SpecialPrice, newPrice.SpecialPrice)

	// Выборка по ID
	model, err = sqlRepo.SelectByID(ctx, model.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.Price, newPrice.Price)
	assert.Equal(t, model.Discount, newPrice.Discount)
	assert.Equal(t, model.SpecialPrice, newPrice.SpecialPrice)

	// Выборка по названию
	models, err := sqlRepo.SelectByCardID(ctx, modelCard.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, model := range models {
		assert.Equal(t, model.Price, newPrice.Price)
		assert.Equal(t, model.Discount, newPrice.Discount)
		assert.Equal(t, model.SpecialPrice, newPrice.SpecialPrice)
	}

	// Обновление
	newPrice.Price = 6.6
	newPrice.Discount = 1.6
	newPrice.SpecialPrice = 9.0
  newPrice.ID = model.ID
  newPrice.CardID = modelCard.ID
  newPrice.SellerID = modelSeller.ID

	err = sqlRepo.UpdateExecOne(ctx, newPrice)
	if err != nil {
		t.Fatal(err)
	}

	// Выборка по ID
	model, err = sqlRepo.SelectByID(ctx, newPrice.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.Price, newPrice.Price)
	assert.Equal(t, model.Discount, newPrice.Discount)
	assert.Equal(t, model.SpecialPrice, newPrice.SpecialPrice)
}
