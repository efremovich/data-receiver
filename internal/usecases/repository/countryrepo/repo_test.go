package countryrepo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/countryrepo"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

func TestCountryRepo(t *testing.T) {
	ctx := context.Background()

	conn, _, err := postgresdb.GetMockConn("../../../../migrations/data_base")
	if err != nil {
		t.Fatalf("ошибка создания мокового соединения. %s", err.Error())
	}

	countryRepo, err := countryrepo.NewCountryRepo(ctx, conn)
	if err != nil {
		t.Fatal(err)
	}

	newCountry := entity.Country{
		Name: uuid.NewString(),
	}

	modelCountry, err := countryRepo.Insert(ctx, newCountry)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, modelCountry.Name, newCountry.Name)

	modelSelectCountry, err := countryRepo.SelectByID(ctx, modelCountry.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, modelCountry.Name, modelSelectCountry.Name)

	modelSelectCountry, err = countryRepo.SelectByName(ctx, modelCountry.Name)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, modelCountry.Name, modelSelectCountry.Name)

	err = countryRepo.UpdateExecOne(ctx, *modelSelectCountry)
	if err != nil {
		t.Fatal(err)
	}
}
