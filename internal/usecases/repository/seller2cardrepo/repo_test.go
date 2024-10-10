package seller2cardrepo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/efremovich/data-receiver/pkg/postgresdb"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/brandrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/cardrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/seller2cardrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/sellerrepo"
)

func TestConvertToDBWb2Card(t *testing.T) {
	ctx := context.Background()

	conn, _, err := postgresdb.GetMockConn("../../../../migrations/data_base")
	if err != nil {
		t.Fatalf("ошибка создания мокового соединения. %s", err.Error())
	}

	// Создание Seller
	sqlSellerRepo, err := sellerrepo.NewSellerRepo(ctx, conn)
	if err != nil {
		t.Fatal(err)
	}

	newSeller := entity.Seller{
		Title:      uuid.NewString(),
		IsEnabled:  true,
		ExternalID: uuid.NewString(),
	}

	modelSeller, err := sqlSellerRepo.Insert(ctx, newSeller)
	if err != nil {
		t.Fatal(err)
	}

	// Создание Brand
	sqlRepo, err := brandrepo.NewBrandRepo(ctx, conn)
	if err != nil {
		t.Fatal(err.Error())
	}

	newBrand := entity.Brand{
		Title:    uuid.NewString(),
		SellerID: modelSeller.ID,
	}

	modelBrand, err := sqlRepo.Insert(ctx, newBrand)
	if err != nil {
		t.Fatal(err)
	}

	// Создание Card
	sqlCardRepo, err := cardrepo.NewCardRepo(ctx, conn)
	if err != nil {
		t.Fatal(err.Error())
	}

	newCard := entity.Card{
		ExternalID:  0,
		VendorID:    uuid.NewString(),
		VendorCode:  uuid.NewString(),
		Title:       uuid.NewString(),
		Description: uuid.NewString(),
		Brand:       *modelBrand,
	}

	modelCard, err := sqlCardRepo.Insert(ctx, newCard)
	if err != nil {
		t.Fatal(err)
	}

	sqlWb2CardRepo, err := seller2cardrepo.NewWb2CardRepo(ctx, conn)
	if err != nil {
		t.Fatal(err)
	}

	newWb2Card := entity.Seller2Card{
		ExternalID: 111,
		KTID:       222,
		NMUUID:     "FDA",
		CardID:     modelCard.ID,
	}

	modelWb2Card, err := sqlWb2CardRepo.Insert(ctx, newWb2Card)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, modelWb2Card.ExternalID, newWb2Card.ExternalID)
	assert.Equal(t, modelWb2Card.KTID, newWb2Card.KTID)
	assert.Equal(t, modelWb2Card.NMUUID, newWb2Card.NMUUID)
	assert.Equal(t, modelWb2Card.CardID, newWb2Card.CardID)

	modelSelectWb2Card, err := sqlWb2CardRepo.SelectByCardID(ctx, newWb2Card.CardID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, modelSelectWb2Card.ExternalID, newWb2Card.ExternalID)
	assert.Equal(t, modelSelectWb2Card.KTID, newWb2Card.KTID)
	assert.Equal(t, modelSelectWb2Card.NMUUID, newWb2Card.NMUUID)
	assert.Equal(t, modelSelectWb2Card.CardID, newWb2Card.CardID)

	modelSelectWb2Card, err = sqlWb2CardRepo.SelectByExternalID(ctx, newWb2Card.ExternalID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, modelSelectWb2Card.ExternalID, newWb2Card.ExternalID)
	assert.Equal(t, modelSelectWb2Card.KTID, newWb2Card.KTID)
	assert.Equal(t, modelSelectWb2Card.NMUUID, newWb2Card.NMUUID)
	assert.Equal(t, modelSelectWb2Card.CardID, newWb2Card.CardID)

	modelSelectWb2Card, err = sqlWb2CardRepo.SelectByKTID(ctx, newWb2Card.KTID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, modelSelectWb2Card.ExternalID, newWb2Card.ExternalID)
	assert.Equal(t, modelSelectWb2Card.KTID, newWb2Card.KTID)
	assert.Equal(t, modelSelectWb2Card.NMUUID, newWb2Card.NMUUID)
	assert.Equal(t, modelSelectWb2Card.CardID, newWb2Card.CardID)

	modelSelectWb2Card, err = sqlWb2CardRepo.SelectByNMUUID(ctx, newWb2Card.NMUUID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, modelSelectWb2Card.ExternalID, newWb2Card.ExternalID)
	assert.Equal(t, modelSelectWb2Card.KTID, newWb2Card.KTID)
	assert.Equal(t, modelSelectWb2Card.NMUUID, newWb2Card.NMUUID)
	assert.Equal(t, modelSelectWb2Card.CardID, newWb2Card.CardID)

	err = sqlWb2CardRepo.UpdateExecOne(ctx, *modelSelectWb2Card)
	if err != nil {
		t.Fatal(err)
	}
}
