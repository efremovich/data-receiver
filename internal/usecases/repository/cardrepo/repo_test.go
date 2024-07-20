package cardrepo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/brandrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/cardrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/sellerrepo"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

func TestCardRepo(t *testing.T) {
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
	newSeller := entity.Seller{
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

	assert.Equal(t, modelCard.VendorID, newCard.VendorID)
	assert.Equal(t, modelCard.VendorCode, newCard.VendorCode)
	assert.Equal(t, modelCard.Title, newCard.Title)
	assert.Equal(t, modelCard.Description, newCard.Description)
	assert.Equal(t, modelCard.Brand.ID, newCard.Brand.ID)

	modelSelectCard, err := sqlCardRepo.SelectByID(ctx, modelCard.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, modelSelectCard.VendorID, newCard.VendorID)
	assert.Equal(t, modelSelectCard.VendorCode, newCard.VendorCode)
	assert.Equal(t, modelSelectCard.Title, newCard.Title)
	assert.Equal(t, modelSelectCard.Description, newCard.Description)
	assert.Equal(t, modelSelectCard.Brand.ID, newCard.Brand.ID)

	modelSelectCard, err = sqlCardRepo.SelectByTitle(ctx, newCard.Title)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, modelSelectCard.VendorID, newCard.VendorID)
	assert.Equal(t, modelSelectCard.VendorCode, newCard.VendorCode)
	assert.Equal(t, modelSelectCard.Title, newCard.Title)
	assert.Equal(t, modelSelectCard.Description, newCard.Description)
	assert.Equal(t, modelSelectCard.Brand.ID, newCard.Brand.ID)

	modelSelectCard, err = sqlCardRepo.SelectByVendorID(ctx, newCard.VendorID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, modelSelectCard.VendorID, newCard.VendorID)
	assert.Equal(t, modelSelectCard.VendorCode, newCard.VendorCode)
	assert.Equal(t, modelSelectCard.Title, newCard.Title)
	assert.Equal(t, modelSelectCard.Description, newCard.Description)
	assert.Equal(t, modelSelectCard.Brand.ID, newCard.Brand.ID)

	modelSelectCard, err = sqlCardRepo.SelectByTitle(ctx, newCard.Title)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, modelSelectCard.VendorID, newCard.VendorID)
	assert.Equal(t, modelSelectCard.VendorCode, newCard.VendorCode)
	assert.Equal(t, modelSelectCard.Title, newCard.Title)
	assert.Equal(t, modelSelectCard.Description, newCard.Description)
	assert.Equal(t, modelSelectCard.Brand.ID, newCard.Brand.ID)
}
