package entity

import aerrors "github.com/efremovich/data-receiver/pkg/aerror"

func AddUserErrorMessages() {
	aerrors.AppendToUserMessages(userMessages)
}

const (
	// Критические.
	WrongMethodErrorID               aerrors.ErrorId = 1_12_1_001 // 1121001
	MissSendReceiptToErrorID         aerrors.ErrorId = 1_12_1_002 // 1121002
	WrongSendReceiptToErrorID        aerrors.ErrorId = 1_12_1_003 // 1121003
	MissContentDispositionErrorID    aerrors.ErrorId = 1_12_1_004 // 1121004
	WrongContentDispositionToErrorID aerrors.ErrorId = 1_12_1_005 // 1121005
	MissFilenameToErrorID            aerrors.ErrorId = 1_12_1_006 // 1121006
	WrongFileNameErrorID             aerrors.ErrorId = 1_12_1_007 // 1121007
	EmptyBodyErrorID                 aerrors.ErrorId = 1_12_1_008 // 1121008
	ParseValidateCMSErrorID          aerrors.ErrorId = 1_12_1_009 // 1121009
	UnknownOperatorErrorID           aerrors.ErrorId = 1_12_1_010 // 1121010
	DisabledOperatorErrorID          aerrors.ErrorId = 1_12_1_011 // 1121011
	TPNoDirErrorID                   aerrors.ErrorId = 1_12_1_012 // 1121012
	TpNoReceiptErrorID               aerrors.ErrorId = 1_12_1_013 // 1121013
	ReceptAndDescErrorID             aerrors.ErrorId = 1_12_1_014 // 1121014
	ReceptAndDirErrorID              aerrors.ErrorId = 1_12_1_015 // 1121015
	LsNoDescErrorID                  aerrors.ErrorId = 1_12_1_016 // 1121016
	NoDescOrReceiptErrorID           aerrors.ErrorId = 1_12_1_017 // 1121017
	LsInLsErrorID                    aerrors.ErrorId = 1_12_1_018 // 1121018
	DirectoryWrongNameErrorID        aerrors.ErrorId = 1_12_1_019 // 1121019
	FileWrongNameErrorID             aerrors.ErrorId = 1_12_1_020 // 1121020

	// Некритические.
	SelectTPErrorID     aerrors.ErrorId = 1_12_0_001 // 1120001
	InsertTPErrorID     aerrors.ErrorId = 1_12_0_002 // 1120002
	InsertEventErrorID  aerrors.ErrorId = 1_12_0_003 // 1120003
	GetOperatorsMap     aerrors.ErrorId = 1_12_0_004 // 1120004
	SaveStorageErrorID  aerrors.ErrorId = 1_12_0_005 // 1120005
	MakeTRKErrorID      aerrors.ErrorId = 1_12_0_006 // 1120006
	BrokerSendErrorID   aerrors.ErrorId = 1_12_0_007 // 1120007
	OpenTXErrorID       aerrors.ErrorId = 1_12_0_008 // 1120008
	UpdateTPErrorID     aerrors.ErrorId = 1_12_0_009 // 1120009
	CommitTXErrorID     aerrors.ErrorId = 1_12_0_010 // 1120010
	InsertFileStructure aerrors.ErrorId = 1_12_0_011 // 1120011
)

// UserMessages Дефолтные сообщения ошибок для пользователей.
//
//nolint:exhaustive // Непонятная ошибка.
var userMessages = map[aerrors.ErrorId]string{
	WrongMethodErrorID:               "Приём ТП поддерживается только методом POST",
	MissSendReceiptToErrorID:         "Отсутствует обязательный заголовок запроса 'Send-Receipt-To'",
	WrongSendReceiptToErrorID:        "Значение заголовка 'Send-Receipt-To' не является корректным URL-адресом",
	MissContentDispositionErrorID:    "Отсутствует обязательный заголовок запроса 'Content-Disposition'",
	WrongContentDispositionToErrorID: "Заголовок 'Content-Disposition' имеет некорректный формат",
	MissFilenameToErrorID:            "Заголовок 'Content-Disposition' не содержит 'filename'",
	WrongFileNameErrorID:             "Имя файла ТП не соответствует формату '[a-zA-Z]{1,100}).cms'",
	EmptyBodyErrorID:                 "Запрос имеет пустое тело",
	SelectTPErrorID:                  "Ошибка выборки ТП из БД",
	InsertTPErrorID:                  "Ошибка вставки ТП в БД",
	InsertEventErrorID:               "Ошибка вставки евента в БД",
	ParseValidateCMSErrorID:          "CMS-пакет имеет некорректный формат или подпись не прошла базовую верификацию",
	GetOperatorsMap:                  "Ошибка получения списка операторов",
	UnknownOperatorErrorID:           "Не удалось определить оператора-отправителя по сертификату",
	DisabledOperatorErrorID:          "Обмен с оператором-обладателем сертификата отключён",
	SaveStorageErrorID:               "Ошибка сохранения файла в хранилище",
	TPNoDirErrorID:                   "Ошибка проверки файловой структуры ТП: ТП не содержит ни одного файла",
	TpNoReceiptErrorID:               "Ошибка проверки файловой структуры ТП: корневая директория содержит файлы, но не содержит receipts.xml",
	ReceptAndDescErrorID:             "Ошибка проверки файловой структуры ТП: корневая директория содержит файлы receipts.xml и description.xml",
	ReceptAndDirErrorID:              "Ошибка проверки файловой структуры ТП: корневая директория содержит файл receipts.xml и директорию",
	LsNoDescErrorID:                  "Ошибка проверки файловой структуры ТП: директория содержит как минимум 1 файл, но не содержит description.xml",
	NoDescOrReceiptErrorID:           "Ошибка проверки файловой структуры ТП: директория не содержит ни receipts.xm, ни description.xml",
	LsInLsErrorID:                    "Ошибка проверки файловой структуры ТП: присутствует вложенная директория внутри ЛС",
	MakeTRKErrorID:                   "Ошибка создания положительной ТРК",
	BrokerSendErrorID:                "Ошибка отправки сообщения в брокер",
	OpenTXErrorID:                    "Ошибка открытия транзакции",
	UpdateTPErrorID:                  "Ошибка обновления ТП в БД",
	CommitTXErrorID:                  "Ошибка фиксации транзакции",
	InsertFileStructure:              "Ошибка вставки файловой структуры ТП в БД",
	DirectoryWrongNameErrorID:        "Имя директории не соответствует формату [a-z0-9]{32}",
	FileWrongNameErrorID:             "Имя файла не соответствует формату [a-z0-9]{32}.(bin|p7s) или (receipts.xml|description.xml)",
}
