package brokerconsumer

import (
	"context"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
	anats "github.com/efremovich/data-receiver/pkg/anats"
)

type Task struct {
	Cursor    int        `json:"cursor"`
	Seller    string     `json:"seller"`
	UpdatedAt *time.Time `json:"updated_at"`
	Limit     int        `json:"limit"`
}

const PackageCreatorStreamName = "package-sender-stream"

const (
	CardCreatorConsumer = "package-sender-q-inbox"
	SubjectForGetCards  = "package-sender.inbox"
)

const ReceiverStreamName = "package-sender-stream"

const ReceiverSubjectNormalPriority = "package-sender-stream.normal_priority"

type handlerForCreatePackageAndSend func(ctx context.Context, desc entity.PackageDescription, retry int, isLastRetry bool) anats.MessageResultEnum
