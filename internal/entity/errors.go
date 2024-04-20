package entity

import aerrors "github.com/efremovich/data-receiver/pkg/aerror"

func AddUserErrorMessages() {
	aerrors.AppendToUserMessages(userMessages)
}

const (
	// Критические.
	UnmarshalErrorID      aerrors.ErrorId = 1_12_1_002 // 1121002
	PackCMSErrorID        aerrors.ErrorId = 1_12_1_004 // 1121004
	UnknownPkgTypeErrorID aerrors.ErrorId = 1_12_1_005 // 1121005
	EmptyBodyErrorID      aerrors.ErrorId = 1_12_1_008 // 1121008

	// Некритические.
	SelectPkgErrorID     aerrors.ErrorId = 1_12_0_001 // 1120001
	InsertPkgErrorID     aerrors.ErrorId = 1_12_0_002 // 1120002
	InsertEventErrorID   aerrors.ErrorId = 1_12_0_003 // 1120003
	GetStorageErrorID    aerrors.ErrorId = 1_12_0_004 // 1120004
	SaveStorageErrorID   aerrors.ErrorId = 1_12_0_005 // 1120005
	BrokerSendErrorID    aerrors.ErrorId = 1_12_0_007 // 1120007
	OpenTXErrorID        aerrors.ErrorId = 1_12_0_008 // 1120008
	UpdatePkgErrorID     aerrors.ErrorId = 1_12_0_009 // 1120009
	CommitTXErrorID      aerrors.ErrorId = 1_12_0_010 // 1120010
	GetDataFromExSources aerrors.ErrorId = 1_12_0_011 // 1120011
)

// UserMessages Дефолтные сообщения ошибок для пользователей.
//
//nolint:exhaustive // Непонятная ошибка.
var userMessages = map[aerrors.ErrorId]string{
	UnmarshalErrorID:      "Ошибка при десериализации данных",
	PackCMSErrorID:        "Ошибка упаковки данных в CMS",
	UnknownPkgTypeErrorID: "Неизвестный тип пакета",
	EmptyBodyErrorID:      "Запрос имеет пустое тело",

	SelectPkgErrorID:   "Ошибка выборки пакета из БД",
	InsertPkgErrorID:   "Ошибка вставки пакета в БД",
	InsertEventErrorID: "Ошибка вставки события в БД",
	GetStorageErrorID:  "Ошибка получения файла из хранилища",
	SaveStorageErrorID: "Ошибка сохранения файла в хранилище",
	BrokerSendErrorID:  "Ошибка отправки сообщения в брокер",
	OpenTXErrorID:      "Ошибка открытия транзакции",
	UpdatePkgErrorID:   "Ошибка обновления пакета в БД",
	CommitTXErrorID:    "Ошибка фиксации транзакции",
}
