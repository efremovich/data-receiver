package brandrepo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/brandrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/sellerrepo"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

func TestBrandRepo(t *testing.T) {
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

	newSeller := entity.Seller{
		Title:      uuid.NewString(),
		IsEnabled:  true,
		ExternalID: uuid.NewString(),
	}

	modelSeller, err := sqlSellerRepo.Insert(ctx, newSeller)
	if err != nil {
		t.Fatal(err)
	}

	sqlRepo, err := brandrepo.NewBrandRepo(ctx, conn)
	if err != nil {
		t.Fatal(err.Error())
	}

	newBrand := entity.Brand{
		Title:    uuid.NewString(),
		SellerID: modelSeller.ID,
	}
	// Создание
	model, err := sqlRepo.Insert(ctx, newBrand)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.Title, newBrand.Title)
	assert.Equal(t, model.SellerID, newBrand.SellerID)

	// Выборка по ID
	model, err = sqlRepo.SelectByID(ctx, model.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.Title, newBrand.Title)
	assert.Equal(t, model.SellerID, newBrand.SellerID)

	// Выборка по названию
	model, err = sqlRepo.SelectByTitleAndSeller(ctx, newBrand.Title, newBrand.SellerID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.Title, newBrand.Title)
	assert.Equal(t, model.SellerID, newBrand.SellerID)

	// Обновление
	newBrand.Title = uuid.NewString()
	newBrand.SellerID = modelSeller.ID
	newBrand.ID = model.ID

	err = sqlRepo.UpdateExecOne(ctx, newBrand)
	if err != nil {
		t.Fatal(err)
	}

	// Выборка по ID
	model, err = sqlRepo.SelectByID(ctx, newBrand.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.Title, newBrand.Title)
	assert.Equal(t, model.SellerID, newBrand.SellerID)
}
