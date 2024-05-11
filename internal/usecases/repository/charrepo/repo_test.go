package charrepo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/cardrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/charrepo"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

func TestCharRepo(t *testing.T) {
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

	sqlRepo, err := charrepo.NewCharRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}

	newChar := entity.Characteristic{
		Title:  uuid.NewString(),
		Value:  []string{"22", "33", "44", "55", "66", "77", "88", "99"},
		CardID: modelCard.ID,
	}

	// Создание
	model, err := sqlRepo.Insert(ctx, newChar)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, newChar.Title, model.Title)
	assert.Equal(t, newChar.Value, model.Value)

	// Выборка по ID
	model, err = sqlRepo.SelectByID(ctx, model.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, newChar.Title, model.Title)
	assert.Equal(t, newChar.Value, model.Value)

	// Выборка по названию
	models, err := sqlRepo.SelectByCardID(ctx, modelCard.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, model := range models {
		assert.Equal(t, newChar.Title, model.Title)
		assert.Equal(t, newChar.Value, model.Value)
	}

	// Обновление
  newChar.ID = model.ID
	newChar.Title = uuid.NewString()
	newChar.Value = []string{"11", "22", "33", "44", "55", "66", "77", "88", "99"}
	newChar.CardID = modelCard.ID

	err = sqlRepo.UpdateExecOne(ctx, newChar)
	if err != nil {
		t.Fatal(err)
	}

	// Выборка по ID
	model, err = sqlRepo.SelectByID(ctx, newChar.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, newChar.Title, model.Title)
	assert.Equal(t, newChar.Value, model.Value)
}
