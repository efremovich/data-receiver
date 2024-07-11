package usecases

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/webapi/wbfetcher"
	aerror "github.com/efremovich/data-receiver/pkg/aerror"
)

func (s *receiverCoreServiceImpl) ReceiveCards(ctx context.Context, sellerTitle string) aerror.AError {
	// // Найдем маркетплейс в базе
	// seller, err := s.sellerRepo.SelectByTitle(ctx, sellerTitle)
	// if err != nil && !errors.Is(err, sql.ErrNoRows) {
	// 	return aerrors.New(ctx, entity.SelectPkgErrorID, err, "ошибка поиска продавца по наименованию %s в БД: %s ", sellerTitle, err.Error())
	// }

	client := wbfetcher.New(ctx, s.cfg.Seller)

	cards, err := client.GetCards(ctx)
	if err != nil {
		return aerror.New(ctx, entity.GetDataFromExSources, err, "ошибка получение данные из внешнего источника %s в БД: %s ", sellerTitle, err.Error())
	}

	for _, card := range cards {
		_, err = s.cardRepo.Insert(ctx, card)
		if err != nil {
			return aerror.New(ctx, entity.SaveStorageErrorID, err, "Ошибка при сохранении карточки товара %s в БД.", card.Title)
		}
	}
	p := entity.Package{
		ID:        0,
		Type:      "",
		Name:      "",
		SendURL:   "",
		Status:    "",
		ErrorText: "",
		ErrorCode: "",
	}
	err = s.brokerPublisher.SendPackage(ctx, &p, "")
	if err != nil {
		return aerror.New(ctx, entity.BrokerSendErrorID, err, "Ошибка постановки задачи в очередь")
	}
	return nil
}
