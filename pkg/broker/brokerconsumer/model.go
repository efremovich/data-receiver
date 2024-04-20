package brokerconsumer

import (
	"context"

	"github.com/efremovich/data-receiver/internal/entity"
	anats "github.com/efremovich/data-receiver/pkg/anats"
)

type Task struct {
	Cursor int    `json:"cursor"`
	Seller string `json:"seller"`
}

const PackageCreatorStreamName = "dara-receiver_card-stream"

const (
	CardCreatorConsumer = "dara-receiver-card-q-inbox"
	SubjectForGetCards  = "card-getter.inbox"
)

const ReceiverStreamName = "receiver_cards"

const ReceiverSubjectNormalPriority = "receiver_cards.normal_priority"

type handlerForCreatePackageAndSend func(ctx context.Context, desc entity.PackageDescription, retry int, isLastRetry bool) anats.MessageResultEnum
