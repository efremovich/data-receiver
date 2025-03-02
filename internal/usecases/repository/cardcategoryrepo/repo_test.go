package cardcategoryrepo_test

import (
	"context"
	"testing"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/brandrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/cardcategoryrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/cardrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/categoryrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/sellerrepo"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCardCategory(t *testing.T) {
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

	// Создание Card
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

	newCategory := entity.Category{
		ExternalID: 1,
		Title:      uuid.NewString(),
		SellerID:   modelSeller.ID,
	}

	sqlCategoryRepo, err := categoryrepo.NewCategoryRepo(ctx, conn)
	if err != nil {
		t.Fatal(err)
	}

	modelCategory, err := sqlCategoryRepo.Insert(ctx, newCategory)
	if err != nil {
		t.Fatal(err)
	}

	newCardCategory := entity.CardCategory{
		CardID:     modelCard.ID,
		CategoryID: modelCategory.ID,
	}

	sqlCardCategoryRepo, err := cardcategoryrepo.NewCardCategory(ctx, conn)
	if err != nil {
		t.Fatal(err)
	}

	model, err := sqlCardCategoryRepo.Insert(ctx, newCardCategory)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.CardID, modelCard.ID)
	assert.Equal(t, model.CategoryID, modelCategory.ID)

	modelSelect, err := sqlCardCategoryRepo.SelectByID(ctx, model.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, modelSelect.CardID, modelCard.ID)
	assert.Equal(t, modelSelect.CategoryID, modelCategory.ID)
	assert.Equal(t, modelSelect.ID, model.ID)

	modelSelect, err = sqlCardCategoryRepo.SelectByCardIDAndCategoryID(ctx, model.CardID, model.CategoryID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, modelSelect.CardID, modelCard.ID)
	assert.Equal(t, modelSelect.CategoryID, modelCategory.ID)
	assert.Equal(t, modelSelect.ID, model.ID)

	err = sqlCardCategoryRepo.Update(ctx, *model)
	if err != nil {
		t.Fatal(err)
	}
}
