package mediafilerepo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/cardrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/mediafilerepo"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

func TestMediaFileRepo(t *testing.T) {
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

	sqlRepo, err := mediafilerepo.NewMediaFileRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}

	newMediaFile := entity.MediaFile{
		Link:   uuid.NewString(),
		CardID: modelCard.ID,
	}

	// Создание
	model, err := sqlRepo.Insert(ctx, newMediaFile)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, newMediaFile.Link, model.Link)

	// Выборка по названию
	models, err := sqlRepo.SelectByCardID(ctx, modelCard.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, model := range models {
		assert.Equal(t, newMediaFile.Link, model.Link)
	}

	// Обновление
	newMediaFile.ID = model.ID
	newMediaFile.Link = uuid.NewString()
	newMediaFile.CardID = modelCard.ID

	err = sqlRepo.UpdateExecOne(ctx, newMediaFile)
	if err != nil {
		t.Fatal(err)
	}

	// Выборка по ID
	models, err = sqlRepo.SelectByCardID(ctx, modelCard.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, model := range models {
		assert.Equal(t, newMediaFile.Link, model.Link)
	}
}
