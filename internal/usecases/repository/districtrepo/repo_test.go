package districtrepo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/districtrepo"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

func TestDistrictRepo(t *testing.T) {
	ctx := context.Background()

	conn, _, err := postgresdb.GetMockConn("../../../../migrations/data_base")
	if err != nil {
		t.Fatalf("ошибка создания мокового соединения. %s", err.Error())
	}

	districtRepo, err := districtrepo.NewDistrictRepo(ctx, conn)
	if err != nil {
		t.Fatal(err)
	}

	newDistrict := entity.District{
		Name: uuid.NewString(),
	}

	modelDistrict, err := districtRepo.Insert(ctx, newDistrict)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, modelDistrict.Name, newDistrict.Name)

	modelSelectDistrict, err := districtRepo.SelectByID(ctx, modelDistrict.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, modelDistrict.Name, modelSelectDistrict.Name)

	modelSelectDistrict, err = districtRepo.SelectByName(ctx, modelDistrict.Name)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, modelDistrict.Name, modelSelectDistrict.Name)

	err = districtRepo.UpdateExecOne(ctx, *modelSelectDistrict)
	if err != nil {
		t.Fatal(err)
	}
}
