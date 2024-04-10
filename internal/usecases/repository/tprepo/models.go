package tprepo

import (
	"context"
	"database/sql"
	"time"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/internal/usecases/repository"
)

type transportPackageDB struct {
	ID         int64          `db:"id"`
	Name       string         `db:"name"`
	IsReceipt  sql.NullBool   `db:"is_receipt"`
	ReceiptURL string         `db:"receipt_url"`
	StatusID   int            `db:"tp_status_id"`
	StatusDesc string         `db:"tp_status_desc"`
	Origin     sql.NullString `db:"sender_operator_code"`

	CreatedAt time.Time `db:"created_at"`

	ErrorCode sql.NullString `db:"error_code"`
	ErrorText sql.NullString `db:"error_text"`
}

func convertToDBTransportPackage(_ context.Context, in entity.TransportPackage, statusesMap map[entity.TpStatusEnum]int) *transportPackageDB {
	return &transportPackageDB{
		ID:         in.ID,
		Name:       in.Name,
		ReceiptURL: in.ReceiptURL,
		IsReceipt:  repository.BoolToNullBoolean(in.IsReceipt),
		StatusID:   statusesMap[in.Status],
		StatusDesc: string(in.Status),
		Origin:     repository.StringToNullString(in.Origin),
		CreatedAt:  in.CreatedAt,
		ErrorText:  repository.StringToNullString(in.ErrorText),
		ErrorCode:  repository.StringToNullString(in.ErrorCode),
	}
}

func (tp transportPackageDB) ConvertToEntityTransportPackage(_ context.Context) *entity.TransportPackage {
	return &entity.TransportPackage{
		ID:         tp.ID,
		Name:       tp.Name,
		ReceiptURL: tp.ReceiptURL,
		IsReceipt:  repository.NullBooleanToBool(tp.IsReceipt),
		Status:     entity.TpStatusEnum(tp.StatusDesc),
		Origin:     repository.NullStringToString(tp.Origin),
		CreatedAt:  tp.CreatedAt,
		ErrorText:  repository.NullStringToString(tp.ErrorText),
		ErrorCode:  repository.NullStringToString(tp.ErrorCode),
	}
}

var tpStatusEnumList = []entity.TpStatusEnum{
	entity.TpStatusEnumNew,
	entity.TpStatusEnumSuccess,
	entity.TpStatusEnumFailed,
	entity.TpStatusEnumFailedInternal,
}

var tpEventTypeEnumList = []entity.TpEventTypeEnum{
	entity.CreatedTpEventType,
	entity.SuccessEventType,
	entity.GotAgainEventType,
	entity.ReprocessEventType,
	entity.ErrorEventType,
	entity.SendTaskNext,
}
