package cardrepo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository/brandrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/cardrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/categoryrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/charrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/pricerepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/sellerrepo"
	"github.com/efremovich/data-receiver/internal/usecases/repository/sizerepo"
	"github.com/efremovich/data-receiver/pkg/postgresdb"
)

func TestCardRepo(t *testing.T) {
	ctx := context.Background()

	conn, _, err := postgresdb.GetMockConn("../../../../migrations/data_receiver_db")
	if err != nil {
		t.Fatalf("ошибка создания мокового соединения. %s", err.Error())
	}

	sqlCardRepo, err := cardrepo.NewCardRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}
	sqlCharRepo, err := charrepo.NewCharRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}
	sqlSizeRepo, err := sizerepo.NewSizeRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}
	sqlPriceRepo, err := pricerepo.NewPriceRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}
	sqlCategoryRepo, err := categoryrepo.NewCategoryRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}
	sqlSellerRepo, err := sellerrepo.NewSellerRepo(ctx, conn)
	if err != nil {
		t.Fatalf(err.Error())
	}

	sqlBrandRepo, err := brandrepo.NewBrandRepo(ctx, conn)
	if err != nil {
		t.Fatal(err)
	}

	err = sqlCardRepo.Ping(context.Background())
	if err != nil {
		t.Fatalf(err.Error())
	}

	newSeller := entity.Seller{
		Title:    "WB",
		IsEnable: true,
		ExtID:    "106-456-456",
	}

	sellerModel, err := sqlSellerRepo.Insert(ctx, newSeller)
	if err != nil {
		t.Fatalf(err.Error())
	}

	assert.Equal(t, newSeller.Title, sellerModel.Title)
	assert.Equal(t, newSeller.IsEnable, sellerModel.IsEnable)
	assert.Equal(t, newSeller.ExtID, sellerModel.ExtID)
	assert.Equal(t, newSeller.ID, sellerModel.ID)

	newBrand := entity.Brand{
		Title:    "Loretto",
		SellerID: sellerModel.ID,
	}

	brandModel, err := sqlBrandRepo.Insert(ctx, newBrand)
	if err != nil {
		t.Fatal(err)
	}

  assert.Equal(t, newBrand.Title, brandModel.Title)
  //TODO Доделать тест
	newCard := entity.Card{
		VendorID:    uuid.NewString(),
		VendorCode:  uuid.NewString(),
		Title:       uuid.NewString(),
		Description: uuid.NewString(),
	}

	cardModel, err := sqlCardRepo.Insert(ctx, newCard)
	if err != nil {
		t.Fatal(err.Error())
	}

	assert.Equal(t, newCard.VendorCode, cardModel.VendorCode)
	assert.Equal(t, newCard.VendorID, cardModel.VendorID)
	assert.Equal(t, newCard.Title, cardModel.Title)
	assert.Equal(t, newCard.Description, cardModel.Description)

	characteristics := []entity.Characteristic{
		{
			Title:  "Цвет",
			Value:  []string{"белый", "черный", "красный"},
			CardID: cardModel.ID,
		},
		{
			Title:  "Тип",
			Value:  []string{"компьютер", "смартфон", "ноутбук"},
			CardID: cardModel.ID,
		},
	}

	for _, char := range characteristics {
		charModel, err := sqlCharRepo.Insert(ctx, char)
		if err != nil {
			t.Fatal(err.Error())
		}
		assert.Equal(t, char.Title, charModel.Title)
		assert.Equal(t, char.Value, charModel.Value)
		assert.Equal(t, char.CardID, charModel.CardID)
	}

	newPrice := entity.Price{
		Price:        5.55,
		Discount:     0.5,
		SpecialPrice: 10.0,
		SellerID:     1,
		CardID:       cardModel.ID,
	}

	priceModel, err := sqlPriceRepo.Insert(ctx, newPrice)
	if err != nil {
		t.Fatal(err.Error())
	}
	assert.Equal(t, newPrice.Price, priceModel.Price)
	assert.Equal(t, newPrice.Discount, priceModel.Discount)
	assert.Equal(t, newPrice.SpecialPrice, priceModel.SpecialPrice)
	assert.Equal(t, newPrice.SellerID, priceModel.SellerID)
	assert.Equal(t, newPrice.CardID, priceModel.CardID)

	sizes := []*entity.Size{
		{
			TechSize: "30-41",
			Title:    "XXL",
			PriceID:  priceModel.ID,
			CardID:   cardModel.ID,
		},
		{
			TechSize: "188",
			Title:    "Рост 188",
			PriceID:  priceModel.ID,
			CardID:   cardModel.ID,
		},
	}

	for _, elem := range sizes {
		model, err := sqlSizeRepo.Insert(ctx, *elem)
		if err != nil {
			t.Fatal(err.Error())
		}
		assert.Equal(t, elem.Title, model.Title)
		assert.Equal(t, elem.TechSize, model.TechSize)
		assert.Equal(t, elem.CardID, model.CardID)
		assert.Equal(t, elem.PriceID, model.PriceID)
	}

	cardModel, err = sqlCardRepo.SelectByID(ctx, cardModel.ID)
	if err != nil {
		t.Fatalf(err.Error())
	}

	assert.Equal(t, newCard.Title, cardModel.Title)

	newCategory := entity.Category{
		Title:    "Шубы",
		CardID:   cardModel.ID,
		SellerID: newSeller.ID,
	}

	categoryModel, err := sqlCategoryRepo.Insert(ctx, newCategory)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, newCategory.Title, categoryModel.Title)
	assert.Equal(t, newCategory.CardID, categoryModel.CardID)
	assert.Equal(t, newCategory.SellerID, categoryModel.SellerID)
}
