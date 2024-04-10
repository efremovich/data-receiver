package usecases

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/efremovich/data-receiver/internal/entity"
	aerror "github.com/efremovich/data-receiver/pkg/aerror"
	"github.com/efremovich/data-receiver/pkg/logger"
)

func (s *receiverCoreServiceImpl) ReceivePackage(ctx context.Context, tpName string, tpBytes []byte, receiptURL string) ([]byte, aerror.AError) {
	tp, err := s.tpRepo.SelectByName(ctx, tpName)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, nil 
	}

	event := entity.GotAgainEventType
	if tp == nil {
		event = entity.CreatedTpEventType

		tp, err = s.tpRepo.Insert(ctx, tpName, receiptURL)
		if err != nil {
			return nil, aerror.New(ctx, entity.InsertTPErrorID, err, "ошибка вставки ТП в БД: %s", err.Error())
		}
	}

	logger.GetLoggerFromContext(ctx).SetAttr("tp_id", tp.ID).Debugf("после select/insert tp")

	err = s.tpRepo.AddNewEvent(tp.ID, event, tp.Name)
	if err != nil {
		return nil, aerror.New(ctx, entity.InsertEventErrorID, err, "ошибка вставки евента в БД: %s", err.Error())
	}

	// Если ТП находится в статусе критической ошибки - вернём ту же ТРК.
	if tp.Status == entity.TpStatusEnumFailed {
		errorID, err := strconv.Atoi(tp.ErrorCode)
		if err != nil {
			// Если в БД по какой-то причине некорректное ID ошибки, то придется переобрабатывать ТП, ничего не поделать.
			logger.GetLoggerFromContext(ctx).Errorf("не удалось запарсить код ошибки тп в инт: %s", tp.ErrorCode)
		} else {
			return nil, aerror.NewCritical(ctx, aerror.ErrorId(errorID), nil, "тп в статусе критической ошибки пришёл повторно")
		}
	}

	// Если ТП уже был успешно обработан - вернём положительную ТРК.
	if tp.Status == entity.TpStatusEnumSuccess {
		return s.makeSuccessTRK(ctx, tp)
	}

	// Основная бизнес-логика в этой функции: проверка серта, проверка файловой структуры, сохранение файлов в сторадж.
	trk, tpFileStructure, aerr := s.processReceivePackage(ctx, tp, tpBytes)
	if aerr != nil {
		// Если в ходе основной бизнес-логики произошла ошибка - сохраним её в БД.
		return nil, s.saveResultToDB(ctx, *tp, nil, aerr)
	}

	// Сохраним результат положительной обработки в БД.
	aerr = s.saveResultToDB(ctx, *tp, tpFileStructure, nil)
	if aerr != nil {
		return nil, aerr
	}

	// Сохраним евент, который подтверждает отправку на следующий этап обработки.
	err = s.tpRepo.AddNewEvent(tp.ID, entity.SendTaskNext, tp.Name)
	if err != nil {
		s.metricsCollector.AddReceiveTPInternalError(entity.InsertEventErrorID.UserMessage())
		logger.GetLoggerFromContext(ctx).Errorf("ошибка сохранения евента о отправке сообщения в брокер в БД: %s", err.Error())
	}

	return trk, nil
}

func (s *receiverCoreServiceImpl) processReceivePackage(ctx context.Context, tp *entity.TransportPackage, tpBytes []byte) ([]byte, []*entity.TpDirectory, aerror.AError) {
	// Эта функция проводит базовую валидацию ТП и разбирает его на сертификат и контент.

	// Получили отпечаток сертификата.
	thumb := ""

	// Получили список операторов, с которыми поддерживаем обмен.
	operators, err := s.operatorRepo.GetOperatorsMap(ctx)
	if err != nil {
		return nil, nil, aerror.New(ctx, entity.GetOperatorsMap, err, "ошибка получения списка операторов: %s", err.Error())
	}

	// По отпечатку нашли оператора-отправителя.
	operator, ok := operators[thumb]
	if !ok {
		return nil, nil, aerror.NewCritical(ctx, entity.UnknownOperatorErrorID, nil, "сертификат %s не принадлежит ни одному из операторов", thumb)
	}

	if operator.IsDisabled {
		return nil, nil, aerror.NewCritical(ctx, entity.DisabledOperatorErrorID, nil, "сертификат %s принадлежит оператору %s, роуминг с которым отключен", thumb, operator.Code)
	}

	tp.Origin = operator.Code

	s.metricsCollector.IncOperator(operator.Name)
	logger.GetLoggerFromContext(ctx).SetAttr("origin", operator.Code).Debugf("прошли основные проверки тп. определили по серту отправителя ТП")

	filecatalog := make(map[string][]byte)
	// Преобразуем каталог файлов в более удобную структуру.
	tpFileStructure := buildFileStructure(ctx, filecatalog)

	// Базово провалидировали файловую структуру ТП.
	aerr := validateFilesStructure(ctx, tpFileStructure)
	if aerr != nil {
		return nil, nil, aerr
	}

	// Если файловая структура валидна - определим, содержит ли ТП квитанцию.
	_, isReceipt := filecatalog[""]
	tp.IsReceipt = &isReceipt

	// Сохраним ТП и его файлы в хранилище.
	err = s.saveToStorage(ctx, tp.Name, tpBytes, tpFileStructure)
	if err != nil {
		return nil, nil, aerror.New(ctx, entity.SaveStorageErrorID, err, "%s", err.Error())
	}

	trk, aerr := s.makeSuccessTRK(ctx, tp)
	if aerr != nil {
		return nil, nil, aerr
	}

	return trk, tpFileStructure, nil
}

func buildFileStructure(_ context.Context, files map[string][]byte) []*entity.TpDirectory {
	directories := map[string]*entity.TpDirectory{}

	for name, bytes := range files {
		// Получаем из пути файла директорию и имя файла.
		fileName := filepath.Base(name)
		dirPath := filepath.Dir(name)

		targetDir, ok := directories[dirPath]
		if !ok {
			targetDir = &entity.TpDirectory{
				Name:  dirPath,
				Files: make(map[string][]byte),
			}

			directories[dirPath] = targetDir
		}

		targetDir.Files[fileName] = bytes
	}

	res := make([]*entity.TpDirectory, 0, len(directories))

	for _, v := range directories {
		res = append(res, v)
	}

	return res
}

// валидация файловой структуры ТП.
func validateFilesStructure(ctx context.Context, tpFileStructure []*entity.TpDirectory) aerror.AError {
	if len(tpFileStructure) == 0 {
		return aerror.NewCritical(ctx, entity.TPNoDirErrorID, nil, "ТП не содержит ни одной директории с файлами")
	}

	// Проверим, что нет вложенности больше 1. То есть dir.Name = 'abc/abc' - некорректно.
	for _, dir := range tpFileStructure {
		if strings.ContainsRune(dir.Name, filepath.Separator) {
			return aerror.NewCritical(ctx, entity.LsInLsErrorID, nil, "ТП содержит вложенные директории: %s", dir.Name)
		}

		// Название директории или '.' (корневая) или должно соответствовать регулярке.
		if dir.Name != "." && !entity.PxDirName.MatchString(dir.Name) {
			return aerror.NewCritical(ctx, entity.DirectoryWrongNameErrorID, nil, "ТП содержит директорию с некорректным названием: %s", dir.Name)
		}

	}

	// Сначала проверим кейс, что это квитанция.
	for _, dir := range tpFileStructure {
		// если корневая директория содержит файлы, значит это должно быть ТП с квитанцией
		if dir.Name == "." && len(dir.Files) > 0 {
			return nil
		}
	}

	return nil
}

func (s *receiverCoreServiceImpl) makeSuccessTRK(ctx context.Context, tp *entity.TransportPackage) ([]byte, aerror.AError) {
	return nil, nil 
}

// Cохраняем ТП и файлы в сторадж.
func (s *receiverCoreServiceImpl) saveToStorage(ctx context.Context, tpName string, tpByes []byte, tpFileStructure []*entity.TpDirectory) error {
	err := s.storage.SaveFile(ctx, tpName, tpByes)
	if err != nil {
		return fmt.Errorf("ошибка сохранения ТП в сторадж: %s", err.Error())
	}

	for _, dir := range tpFileStructure {
		for filename, filedata := range dir.Files {
			err = s.storage.SaveFile(ctx, filename, filedata)
			if err != nil {
				return fmt.Errorf("ошибка сохранения файла %s в сторадж: %s", filename, err.Error())
			}
		}
	}

	return nil
}

func (s *receiverCoreServiceImpl) saveResultToDB(ctx context.Context, tp entity.TransportPackage, tpFileStructure []*entity.TpDirectory, res aerror.AError) aerror.AError {
	var (
		event     entity.TpEventTypeEnum
		eventDesc string
	)

	if res == nil {
		tp.ErrorCode = ""
		tp.ErrorText = ""
		event = entity.SuccessEventType
		eventDesc = tp.Name
		tp.Status = entity.TpStatusEnumSuccess
	} else {
		tp.ErrorCode = res.Code()
		tp.ErrorText = res.DeveloperMessage()
		event = entity.ErrorEventType
		eventDesc = res.Code()

		if res.IsCritical() {
			tp.Status = entity.TpStatusEnumFailed
		} else {
			tp.Status = entity.TpStatusEnumFailedInternal
		}
	}

	tx, err := s.tpRepo.BeginTX(ctx)
	if err != nil {
		return aerror.New(ctx, entity.OpenTXErrorID, err, "ошибка при создании транзакции: %s", err.Error())
	}

	defer func() { _ = tx.RollbackIfNotCommitted() }()

	err = s.tpRepo.WithTx(&tx).AddNewEvent(tp.ID, event, eventDesc)
	if err != nil {
		return aerror.New(ctx, entity.InsertEventErrorID, err, "ошибка при вставке евента в БД: %s", err.Error())
	}

	err = s.tpRepo.WithTx(&tx).UpdateExecOne(ctx, tp)
	if err != nil {
		return aerror.New(ctx, entity.UpdateTPErrorID, err, "ошибка при обновлении ТП в БД: %s", err.Error())
	}

	if res == nil && tpFileStructure != nil {
		err = s.tpRepo.WithTx(&tx).SaveFileStructure(ctx, tp.ID, tpFileStructure)
		if err != nil {
			return aerror.New(ctx, entity.InsertFileStructure, err, "ошибка при сохранении файловой структуры ТП в БД: %s", err.Error())
		}
	}

	if err = tx.Commit(); err != nil {
		return aerror.New(ctx, entity.CommitTXErrorID, err, "ошибка при коммите транзакции: %s", err.Error())
	}

	return res
}
