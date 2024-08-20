package regionrepo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/countryrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/districtrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/regionrepo"
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

	regionRepo, err := regionrepo.NewRegionRepo(ctx, conn)
	if err != nil {
		t.Fatal(err)
	}

	newRegion := entity.Region{
		RegionName: uuid.NewString(),
		District:   entity.District{ID: modelDistrict.ID},
		Country:    entity.Country{ID: modelCountry.ID},
	}

	modelRegion, err := regionRepo.Insert(ctx, &newRegion)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, modelRegion.Country.ID, newRegion.Country.ID)
	assert.Equal(t, modelRegion.RegionName, newRegion.RegionName)
	assert.Equal(t, modelRegion.District.ID, newRegion.District.ID)

	modelSelectRegion, err := regionRepo.SelectByName(ctx, modelRegion.RegionName, modelRegion.Country.ID, modelRegion.District.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, modelRegion.Country.ID, newRegion.Country.ID)
	assert.Equal(t, modelRegion.RegionName, newRegion.RegionName)
	assert.Equal(t, modelRegion.District.ID, newRegion.District.ID)

	err = regionRepo.UpdateExecOne(ctx, modelSelectRegion)
	if err != nil {
		t.Fatal(err)
	}
}
