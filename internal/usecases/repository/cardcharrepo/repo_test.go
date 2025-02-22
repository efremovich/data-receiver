package cardcharrepo_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/brandrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/cardcharrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/cardrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/charrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/sellerrepo"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

func TestCharCardRepo(t *testing.T) {
	ctx := context.Background()

	conn, _, err := postgresdb.GetMockConn("../../../../migrations/data_base")
	if err != nil {
		t.Fatalf("ошибка создания мокового соединения. %s", err.Error())
	}

	// Создание Seller
	sellerRepo, err := sellerrepo.NewSellerRepo(ctx, conn)
	if err != nil {
		t.Fatal(err)
	}

	newSeller := entity.MarketPlace{
		Title:      uuid.NewString(),
		IsEnabled:  true,
		ExternalID: uuid.NewString(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	modelSeller, err := sellerRepo.Insert(ctx, newSeller)
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

	// Создание Card
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
	// Создание Characteristic
	sqlCharacteristicRepo, err := charrepo.NewCharRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}

	newChar := entity.Characteristic{
		Title: uuid.NewString(),
	}

	modelChar, err := sqlCharacteristicRepo.Insert(ctx, newChar)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// Создание CardCharacteristic
	sqlCardCharacteristicRepo, err := cardcharrepo.NewCharRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}
	newCardChar := entity.CardCharacteristic{
		Value:            []string{"1", "2", "3"},
		Title:            modelChar.Title,
		CharacteristicID: modelChar.ID,
		CardID:           modelCard.ID,
	}

	modelCardChar, err := sqlCardCharacteristicRepo.Insert(ctx, newCardChar)
	if err != nil {
		t.Fatalf(err.Error())
	}
	assert.Equal(t, modelCardChar.Title, newCardChar.Title)
	assert.Equal(t, modelCardChar.CharacteristicID, newCardChar.CharacteristicID)
	assert.Equal(t, modelCardChar.CardID, newCardChar.CardID)
	assert.Equal(t, modelCardChar.Value, newCardChar.Value)

	modelSelectCardChar, err := sqlCardCharacteristicRepo.SelectByID(ctx, modelCardChar.CardID)
	if err != nil {
		t.Fatalf(err.Error())
	}
	assert.Equal(t, modelCardChar.Title, modelSelectCardChar.Title)
	assert.Equal(t, modelCardChar.CharacteristicID, modelSelectCardChar.CharacteristicID)
	assert.Equal(t, modelCardChar.CardID, modelSelectCardChar.CardID)
	assert.Equal(t, modelCardChar.Value, modelSelectCardChar.Value)

	modelSelectCardChar, err = sqlCardCharacteristicRepo.SelectByCardIDAndCharID(ctx, modelCard.ID, modelChar.ID)
	if err != nil {
		t.Fatalf(err.Error())
	}
	assert.Equal(t, modelCardChar.Title, modelSelectCardChar.Title)
	assert.Equal(t, modelCardChar.CharacteristicID, modelSelectCardChar.CharacteristicID)
	assert.Equal(t, modelCardChar.CardID, modelSelectCardChar.CardID)
	assert.Equal(t, modelCardChar.Value, modelSelectCardChar.Value)
}
