package sizerepo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/sizerepo"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

func TestSizeRepo(t *testing.T) {
	ctx := context.Background()

	conn, _, err := postgresdb.GetMockConn("../../../../migrations/data_base")
	if err != nil {
		t.Fatalf("ошибка создания мокового соединения. %s", err.Error())
	}
  // Создание размера
	sqlRepo, err := sizerepo.NewSizeRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}
	newSize := entity.Size{
		TechSize: uuid.NewString(),
		Title:    uuid.NewString(),
	}
	// Создание
	model, err := sqlRepo.Insert(ctx, newSize)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.TechSize, newSize.TechSize)
	assert.Equal(t, model.Title, newSize.Title)

	// Выборка по ID
	model, err = sqlRepo.SelectByID(ctx, model.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.TechSize, newSize.TechSize)
	assert.Equal(t, model.Title, newSize.Title)

	// Выборка по Title
	model, err = sqlRepo.SelectByTitle(ctx, model.Title)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.TechSize, newSize.TechSize)
	assert.Equal(t, model.Title, newSize.Title)

	// Выборка по Title
	model, err = sqlRepo.SelectByTechSize(ctx, model.TechSize)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.TechSize, newSize.TechSize)
	assert.Equal(t, model.Title, newSize.Title)

	// Обновление
	newSize.Title = uuid.NewString()
	newSize.TechSize = uuid.NewString()
	newSize.ID = model.ID

	err = sqlRepo.UpdateExecOne(ctx, newSize)
	if err != nil {
		t.Fatal(err)
	}

	// Выборка по ID
	model, err = sqlRepo.SelectByID(ctx, newSize.ID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, model.TechSize, newSize.TechSize)
	assert.Equal(t, model.Title, newSize.Title)
}
