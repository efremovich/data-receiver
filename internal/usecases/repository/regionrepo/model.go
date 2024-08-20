package regionrepo

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
)

type regionDB struct {
	ID         int64  `db:"id"`
	RegionName string `db:"region_name"`
	DistrictID int64  `db:"district_id"`
	CountryID  int64  `db:"country_id"`
}

func convertToDBRegion(_ context.Context, in *entity.Region) *regionDB {
	return &regionDB{
		ID:         in.ID,
		RegionName: in.RegionName,
		DistrictID: in.District.ID,
		CountryID:  in.Country.ID,
	}
}

func (r regionDB) convertToEntityRegion(_ context.Context) *entity.Region {
	return &entity.Region{
		ID:         r.ID,
		RegionName: r.RegionName,
		District:   entity.District{ID: r.DistrictID},
		Country:    entity.Country{ID: r.CountryID},
	}
}
