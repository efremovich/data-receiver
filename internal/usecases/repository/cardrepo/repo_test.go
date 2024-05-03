package cardrepo_test

import (
	"context"
	"testing"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/cardrepo"
	"github.com/efremovich/data-receiver/pkg/postgresdb"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCardRepo(t *testing.T) {
	ctx := context.Background()

	conn, _, err := postgresdb.GetMockConn("../../../../migrations/data_receiver_db")
	if err != nil {
		t.Fatalf("ошибка создания мокового соединения. %s", err.Error())
	}

	sqlCardRepo, err := cardrepo.NewCardRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}
	mockRepo := cardrepo.NewCardRepoMock()

	repos := []cardrepo.CardRepo{sqlCardRepo, mockRepo}

	for _, cardRepo := range repos {
		err = cardRepo.Ping(context.Background())
		if err != nil {
			t.Fatalf(err.Error())
		}
		newCard := entity.Card{
			VendorID:    uuid.NewString(),
			VendorCode:  uuid.NewString(),
			Title:       uuid.NewString(),
			Description: uuid.NewString(),
		}

		cardModel, err := cardRepo.Insert(ctx, newCard)
		if err != nil {
			t.Fatal(err.Error())
		}

		assert.Equal(t, newCard.VendorCode, cardModel.VendorCode)
		assert.Equal(t, newCard.VendorID, cardModel.VendorID)
		assert.Equal(t, newCard.Title, cardModel.Title)
		assert.Equal(t, newCard.Description, cardModel.Description)

	}
}
