package brandrepo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/brandrepo"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

func TestBrandRepo(t *testing.T) {
	ctx := context.Background()

	conn, _, err := postgresdb.GetMockConn("../../../../migrations/data_receiver_db")
	if err != nil {
		t.Fatalf("ошибка создания мокового соединения. %s", err.Error())
	}

	sqlRepo, err := brandrepo.NewBrandRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}

	newBrand := entity.Brand{
		Title:    uuid.NewString(),
		SellerID: 1,
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
	model, err = sqlRepo.SelectByTitle(ctx, newBrand.Title)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.Title, newBrand.Title)
	assert.Equal(t, model.SellerID, newBrand.SellerID)

	// Обновление
	newBrand.Title = uuid.NewString()
	newBrand.SellerID = 2
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