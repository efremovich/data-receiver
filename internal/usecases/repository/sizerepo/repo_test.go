package sizerepo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/cardrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/pricerepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/sellerrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/sizerepo"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

func TestSizeRepo(t *testing.T) {
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
		ExternalID:    uuid.NewString(),
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

	sqlRepo, err := sizerepo.NewSizeRepo(ctx, conn)
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
	model, err := sqlRepo.Insert(ctx, newSize)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.TechSize, newSize.TechSize)
	assert.Equal(t, model.Title, newSize.Title)
	assert.Equal(t, model.CardID, newSize.CardID)
	assert.Equal(t, model.PriceID, newSize.PriceID)

	// Выборка по ID
	model, err = sqlRepo.SelectByID(ctx, model.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.TechSize, newSize.TechSize)
	assert.Equal(t, model.Title, newSize.Title)
	assert.Equal(t, model.CardID, newSize.CardID)
	assert.Equal(t, model.PriceID, newSize.PriceID)

	// Выборка по id карточки
	models, err := sqlRepo.SelectByCardID(ctx, newSize.CardID)
	if err != nil {
		t.Fatal(err)
	}
	for _, model := range models {
		assert.Equal(t, model.TechSize, newSize.TechSize)
		assert.Equal(t, model.Title, newSize.Title)
		assert.Equal(t, model.CardID, newSize.CardID)
		assert.Equal(t, model.PriceID, newSize.PriceID)
	}

	// Выборка по id карточки
	models, err = sqlRepo.SelectByPriceID(ctx, newSize.CardID)
	if err != nil {
		t.Fatal(err)
	}
	for _, model := range models {
		assert.Equal(t, model.TechSize, newSize.TechSize)
		assert.Equal(t, model.Title, newSize.Title)
		assert.Equal(t, model.CardID, newSize.CardID)
		assert.Equal(t, model.PriceID, newSize.PriceID)
	}

	// Обновление
	newSize.Title = uuid.NewString()
	newSize.TechSize = uuid.NewString()
	newSize.CardID = modelCard.ID
	newSize.PriceID = modelPrice.ID
	newSize.ID = model.ID

	err = sqlRepo.UpdateExecOne(ctx, newSize)
	if err != nil {
		t.Fatal(err)
	}

	// Выборка по ID
	model, err = sqlRepo.SelectByID(ctx, newSize.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.TechSize, newSize.TechSize)
	assert.Equal(t, model.Title, newSize.Title)
	assert.Equal(t, model.CardID, newSize.CardID)
	assert.Equal(t, model.PriceID, newSize.PriceID)
}
