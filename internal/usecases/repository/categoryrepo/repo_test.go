package categoryrepo_test

import (
	"context"
	"testing"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/categoryrepo"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCategoryRepo(t *testing.T) {
	ctx := context.Background()

	conn, _, err := postgresdb.GetMockConn("../../../../migrations/data_base")
	if err != nil {
		t.Fatalf("ошибка создания мокового соединения. %s", err.Error())
	}

	sqlRepo, err := categoryrepo.NewCategoryRepo(ctx, conn)
	if err != nil {
		t.Fatal(err)
	}

	newCategory := entity.Category{
		ExternalID: 1,
		Title:      uuid.NewString(),
		SellerID:   1,
	}
	// Вставка
	model, err := sqlRepo.Insert(ctx, newCategory)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, model.Title, newCategory.Title)
	assert.Equal(t, model.SellerID, newCategory.Title)
	assert.Equal(t, model.ExternalID, newCategory.ExternalID)

	// Выборка по ID
	model, err = sqlRepo.SelectByID(ctx, model.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.Title, newCategory.Title)
	assert.Equal(t, model.SellerID, newCategory.Title)
	assert.Equal(t, model.ExternalID, newCategory.ExternalID)

	// Выборка по ID
	models, err := sqlRepo.SelectBySellerID(ctx, model.SellerID)
	if err != nil {
		t.Fatal(err)
	}
	for _, model := range models {
		assert.Equal(t, model.Title, newCategory.Title)
		assert.Equal(t, model.SellerID, newCategory.Title)
		assert.Equal(t, model.ExternalID, newCategory.ExternalID)

	}
}
