package mediafilerepo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/brandrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/cardrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/mediafilerepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/mediafiletypeenumrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/sellerrepo"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

func TestMediaFileRepo(t *testing.T) {
	ctx := context.Background()

	conn, _, err := postgresdb.GetMockConn("../../../../migrations/data_base")
	if err != nil {
		t.Fatalf("ошибка создания мокового соединения. %s", err.Error())
	}

	// Создание продавца
	sqlSellerRepo, err := sellerrepo.NewSellerRepo(ctx, conn)
	if err != nil {
		t.Fatal(err)
	}
	newSeller := entity.MarketPlace{
		Title:      uuid.NewString(),
		IsEnabled:  true,
		ExternalID: uuid.NewString(),
	}
	modelSeller, err := sqlSellerRepo.Insert(ctx, newSeller)
	if err != nil {
		t.Fatal(err)
	}
	sqlBrandRepo, err := brandrepo.NewBrandRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}

	newBrand := entity.Brand{
		Title:    uuid.NewString(),
		SellerID: modelSeller.ID,
	}
	// Создание
	modelBrand, err := sqlBrandRepo.Insert(ctx, newBrand)
	if err != nil {
		t.Fatal(err)
	}
	// Создание карточки
	sqlCardRepo, err := cardrepo.NewCardRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}
	newCard := entity.Card{
		Brand:       *modelBrand,
		VendorID:    uuid.NewString(),
		VendorCode:  uuid.NewString(),
		Title:       uuid.NewString(),
		Description: uuid.NewString(),
	}

	modelCard, err := sqlCardRepo.Insert(ctx, newCard)
	if err != nil {
		t.Fatal(err)
	}

	sqlMediaFileTypeEnumRepo, err := mediafiletypeenumrepo.NewMediaFileTypeEnumRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}

	newMediaFileTypeEnum := entity.MediaFileTypeEnum{
		Type: "VIDEO",
	}

	// Создание
	modelMediaFileTypeEnum, err := sqlMediaFileTypeEnumRepo.Insert(ctx, newMediaFileTypeEnum)
	if err != nil {
		t.Fatal(err)
	}

	sqlMediaFileRepo, err := mediafilerepo.NewMediaFileRepo(ctx, conn)
	if err != nil {
		t.Fatal(err)
	}

	newMediaFile := entity.MediaFile{
		Link:   uuid.NewString(),
		CardID: modelCard.ID,
		TypeID: modelMediaFileTypeEnum.ID,
	}

	// Создание
	model, err := sqlMediaFileRepo.Insert(ctx, newMediaFile)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, newMediaFile.Link, model.Link)

	// Выборка по названию
	models, err := sqlMediaFileRepo.SelectByCardID(ctx, modelCard.ID, model.Link)
	if err != nil {
		t.Fatal(err)
	}
	for _, model := range models {
		assert.Equal(t, newMediaFile.Link, model.Link)
	}

	// Обновление
	newMediaFile.ID = model.ID
	newMediaFile.Link = uuid.NewString()
	newMediaFile.CardID = modelCard.ID
	newMediaFile.TypeID = modelMediaFileTypeEnum.ID

	err = sqlMediaFileRepo.UpdateExecOne(ctx, newMediaFile)
	if err != nil {
		t.Fatal(err)
	}
}
