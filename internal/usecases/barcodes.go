package usecases

import (
	"context"
	"errors"

	"github.com/efremovich/data-receiver/internal/entity"
)

func (s *receiverCoreServiceImpl) setBarcode(ctx context.Context, in entity.Barcode) (*entity.Barcode, error) {
	barcode, err := s.barcodeRepo.SelectByBarcode(ctx, in.Barcode)
	if errors.Is(err, ErrObjectNotFound) {
		barcode, err = s.barcodeRepo.Insert(ctx, in)
		if err != nil {
			return nil, err
		}
	}

	return barcode, nil
}

func (s *receiverCoreServiceImpl) getBarcode(ctx context.Context, barcode string) (*entity.Barcode, error) {
	in, err := s.barcodeRepo.SelectByBarcode(ctx, barcode)
	if err != nil {
		return nil, err
	}
	return in, nil
}
