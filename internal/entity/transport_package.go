package entity

import (
	"regexp"
	"time"
)

var (
	// Имя ТП.
	RxTpName = regexp.MustCompile("(?i)^[a-z0-9_-]{1,100}.cms$") // (?i) - регистронезависимый.
	// Имя директории.
	PxDirName = regexp.MustCompile("^[a-z0-9]{32}$")
	// Имя файла.
	PxFileName = regexp.MustCompile(`^[a-z0-9]{32}.(bin|p7s)$`)
)

type TransportPackage struct {
	ID         int64  // id в бд
	Name       string // id транспортного пакета, принятого от внешней системы
	IsReceipt  *bool  // содержит ли ТП ТРК.
	Origin     string // код оператора, отправившего пакет в приемник (определяется по подписи пакета)
	ReceiptURL string // URL, по которому нужно отправить технологическую квитанцию
	Status     TpStatusEnum

	CreatedAt time.Time

	ErrorCode string
	ErrorText string
}

type TpStatusEnum string

const (
	TpStatusEnumNew            TpStatusEnum = "new"
	TpStatusEnumSuccess        TpStatusEnum = "success"
	TpStatusEnumFailed         TpStatusEnum = "failed"
	TpStatusEnumFailedInternal TpStatusEnum = "failed_internal"
)

type TpDirectory struct {
	Name  string            // ID ЛС (по факту - название директории. Для корневной - '.')
	Files map[string][]byte // Карта файлов.
}
