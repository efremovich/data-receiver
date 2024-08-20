package statusrepo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	
  "github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/statusrepo"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)


func TestDistrictRepo(t *testing.T) {
	ctx := context.Background()

	conn, _, err := postgresdb.GetMockConn("../../../../migrations/data_base")
	if err != nil {
		t.Fatalf("ошибка создания мокового соединения. %s", err.Error())
	}

	statusRepo, err := statusrepo.NewStatusRepo(ctx, conn)
	if err != nil {
		t.Fatal(err)
	}

	newStatus := entity.Status{
		Name: uuid.NewString(),
	}

	modelStatus, err := statusRepo.Insert(ctx, newStatus)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, modelStatus.Name, newStatus.Name)


  modelSelectDistrict, err := statusRepo.SelectByName(ctx, modelStatus.Name)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, modelStatus.Name, modelSelectDistrict.Name)

	err = statusRepo.UpdateExecOne(ctx, *modelSelectDistrict)
	if err != nil {
		t.Fatal(err)
	}
}
