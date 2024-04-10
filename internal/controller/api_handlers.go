package controller

import (
	"context"
	"time"

	"github.com/gogo/status"
	"google.golang.org/grpc/codes"

	"github.com/efremovich/data-receiver/internal/entity"

	package_receiver "github.com/efremovich/data-receiver/pkg/package-receiver-service"
)

func (gw *grpcGatewayServerImpl) GetTP(ctx context.Context, in *package_receiver.GetTPRequest) (*package_receiver.GetTPResponse, error) {
	if in.GetDoc() == "" && in.GetTp() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "не переданы параметры для поиска")
	}

	if in.GetDoc() != "" && in.GetTp() != "" {
		return nil, status.Errorf(codes.InvalidArgument, "одновременно переданы параметры doc и tp")
	}


	tp, event, dirs, err := gw.packageReceiver.GetTP(ctx, in.GetTp(), in.GetDoc())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "ошибка при получении ТП из БД: %s", err.Error())
	}

	var result package_receiver.GetTPResponse
	if tp == nil {
		return &result, nil
	}

	result.Founded = true
	result.Tp = tp.Name
	result.Origin = tp.Origin
	result.ReceiptUrl = tp.ReceiptURL
	result.ErrorCode = tp.ErrorCode
	result.ErrorText = tp.ErrorText
	result.CreatedAt = tp.CreatedAt.Format(time.RFC3339)
	result.TimeLayout = time.RFC3339

	if tp.IsReceipt != nil {
		result.IsReceipt = *tp.IsReceipt
	}

	switch tp.Status {
	case entity.TpStatusEnumSuccess:
		result.IsSuccess = true
	case entity.TpStatusEnumFailed:
		result.IsValidationError = true
	case entity.TpStatusEnumFailedInternal:
		result.IsInternalError = true
	case entity.TpStatusEnumNew:
		fallthrough
	default:
		result.IsNew = true
	}

	for _, event := range event {
		if event.EventType == entity.SendTaskNext {
			result.SendTaskNextAt = event.CreatedAt.Format(time.RFC3339)
		}
	}

	for _, dir := range dirs {
		filesNames := make([]string, 0, len(dir.Files))

		for filename := range dir.Files {
			filesNames = append(filesNames, filename)
		}

		result.Content = append(result.GetContent(), &package_receiver.Directory{Name: dir.Name, Files: filesNames})
	}

	return &result, nil
}
