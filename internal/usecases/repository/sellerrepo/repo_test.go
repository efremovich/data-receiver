package sellerrepo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/sellerrepo"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

func TestSellerRepo(t *testing.T) {
	ctx := context.Background()

	conn, _, err := postgresdb.GetMockConn("../../../../migrations/data_receiver_db")
	if err != nil {
		t.Fatalf("ошибка создания мокового соединения. %s", err.Error())
	}

	sqlRepo, err := sellerrepo.NewSellerRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}

	newSeller := entity.Seller{
		Title:    uuid.NewString(),
		ExtID:    uuid.NewString(),
		IsEnable: true,
	}
	// Создание
	model, err := sqlRepo.Insert(ctx, newSeller)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.Title, newSeller.Title)
	assert.Equal(t, model.IsEnable, newSeller.IsEnable)
	assert.Equal(t, model.ExtID, newSeller.ExtID)

	// Выборка по ID
	model, err = sqlRepo.SelectByID(ctx, model.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.Title, newSeller.Title)
	assert.Equal(t, model.IsEnable, newSeller.IsEnable)
	assert.Equal(t, model.ExtID, newSeller.ExtID)

	// Выборка по названию
	model, err = sqlRepo.SelectByTitle(ctx, newSeller.Title)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.Title, newSeller.Title)
	assert.Equal(t, model.IsEnable, newSeller.IsEnable)
	assert.Equal(t, model.ExtID, newSeller.ExtID)

	// Обновление
	newSeller.Title = uuid.NewString()
	newSeller.ExtID = uuid.NewString()
	newSeller.IsEnable = false
	newSeller.ID = model.ID

	err = sqlRepo.UpdateExecOne(ctx, newSeller)
	if err != nil {
		t.Fatal(err)
	}

	// Выборка по ID
	model, err = sqlRepo.SelectByID(ctx, newSeller.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.Title, newSeller.Title)
	assert.Equal(t, model.IsEnable, newSeller.IsEnable)
	assert.Equal(t, model.ExtID, newSeller.ExtID)
}
