package usecases

import (
	"context"
	"errors"

	"github.com/efremovich/data-receiver/internal/entity"
)

func (s *receiverCoreServiceImpl) setCountry(ctx context.Context, in entity.Country) (*entity.Country, error) {
	country, err := s.countryrepo.SelectByName(ctx, in.Name)
	if errors.Is(err, ErrObjectNotFound) {
		country, err = s.countryrepo.Insert(ctx, in)
	}

	if err != nil {
		return nil, err
	}

	return country, nil
}

func (s *receiverCoreServiceImpl) getCountry(ctx context.Context, title string) (*entity.Country, error) {
	country, err := s.countryrepo.SelectByName(ctx, title)
	if err != nil {
		return nil, err
	}
	return country, nil
}

func (s *receiverCoreServiceImpl) getDistrict(ctx context.Context, title string) (*entity.District, error) {
	district, err := s.districtrepo.SelectByName(ctx, title)
	if err != nil {
		return nil, err
	}
	return district, nil
}

func (s *receiverCoreServiceImpl) setDistrict(ctx context.Context, in entity.District) (*entity.District, error) {
	district, err := s.districtrepo.SelectByName(ctx, in.Name)
	if errors.Is(err, ErrObjectNotFound) {
		district, err = s.districtrepo.Insert(ctx, in)
	}

	if err != nil {
		return nil, err
	}

	return district, nil
}

func (s *receiverCoreServiceImpl) getRegion(ctx context.Context, in entity.Region) (*entity.Region, error) {
	region, err := s.regionrepo.SelectByName(ctx, in.RegionName, in.District.ID, in.Country.ID)
	if err != nil {
		return nil, err
	}

	return region, nil
}

func (s *receiverCoreServiceImpl) setRegion(ctx context.Context, in entity.Region) (*entity.Region, error) {
	region, err := s.regionrepo.SelectByName(ctx, in.RegionName, in.District.ID, in.Country.ID)
	if errors.Is(err, ErrObjectNotFound) {
		region, err = s.regionrepo.Insert(ctx, in)
	}

	if err != nil {
		return nil, err
	}

	return region, nil
}
