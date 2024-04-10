package controller

import (
	"context"
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	aerror "github.com/efremovich/data-receiver/pkg/aerror"
	"github.com/efremovich/data-receiver/pkg/alogger"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/pkg/logger"
)

// PackageReceiveV1Handler - хэндлер для обработки ТП.
func (gw *grpcGatewayServerImpl) PackageReceiveV1Handler(req *fiber.Ctx) error {
	gw.metricsCollector.IncCurrentTpInWork()
	gw.metricsCollector.IncTPCounter()

	//nolint: staticcheck // используется строка.
	ctx := context.WithValue(context.Background(), alogger.TraceIdKey, uuid.NewString())
	ctx = logger.AddLoggerInContext(ctx, alogger.NewALogger(ctx, nil))

	var (
		startTime      = time.Now()
		responseBytes  []byte
		responseStatus = http.StatusOK
	)

	defer func() {
		req.Request().ResetBody()
		logger.GetLoggerFromContext(ctx).Flush()
		gw.metricsCollector.DecCurrentTpInWork()
		gw.metricsCollector.AddTPProcessTime(time.Since(startTime))
	}()

	// Провалидировали HTTP-запрос и достали из него имя ТП, ссылку для ответной ТК и байты ТП.
	tpName, receiptURL, tpBytes, aerr := parseAndValidateReq(ctx, req)
	if aerr == nil {
		logger.GetLoggerFromContext(ctx).SetAttr("tp_name", tpName).SetAttr("receipt_url", receiptURL).Infof("валидация запроса прошла успешно. начинается обработка ТП.")

		// Если при валидации не возникло ошибки - тут вызов основной логики обработки ТП.
		responseBytes, aerr = gw.packageReceiver.ReceivePackage(ctx, tpName, tpBytes, receiptURL)
	}

	// В зависимости от ошибки определим статус ответа. Критическая - 400. Не критическая (временная, внутренняя) - 500.
	if aerr != nil {
		if aerr.IsCritical() {
			responseStatus = http.StatusBadRequest
			responseBytes = []byte(aerr.GetID().UserMessage())

			gw.metricsCollector.AddReceiveTPCriticalError(aerr.GetID().UserMessage())
		} else {
			responseStatus = http.StatusInternalServerError
			responseBytes = []byte(fmt.Sprintf("Временная ошибка: код: %s. Попробуйте отправить пакет ещё раз.", aerr.Code()))

			gw.metricsCollector.AddReceiveTPInternalError(aerr.GetID().UserMessage())
		}
	}

	if responseStatus == http.StatusOK {
		logger.GetLoggerFromContext(ctx).Infof("завершена обработка ТП. положительная ТРК. время: %.3fs", time.Since(startTime).Seconds())
	} else {
		logger.GetLoggerFromContext(ctx).Errorf("завершена обработка ТП. отрицательная ТРК. статус %d, текст %s. время: %.3fs. ошибка: %s - %s",
			responseStatus, string(responseBytes), time.Since(startTime).Seconds(), aerr.Code(), aerr.DeveloperMessage())
	}

	req.Response().SetStatusCode(responseStatus)
	req.Response().AppendBody(responseBytes)

	return nil
}

func parseAndValidateReq(ctx context.Context, req *fiber.Ctx) (string, string, []byte, aerror.AError) {
	if req.Method() != http.MethodPost {
		return "", "", nil, aerror.NewCritical(ctx, entity.WrongMethodErrorID, nil, "попытка загрузки методом %s", req.Method())
	}

	receiptURL := string(req.Request().Header.Peek("Send-Receipt-To"))
	if receiptURL == "" {
		return "", "", nil, aerror.NewCritical(ctx, entity.MissSendReceiptToErrorID, nil, "в запросе нет заголовка Send-Receipt-To")
	}

	_, err := url.ParseRequestURI(receiptURL)
	if err != nil {
		return "", "", nil, aerror.NewCritical(ctx, entity.WrongSendReceiptToErrorID, err, "указан некорректный Send-Receipt-To: %s", receiptURL)
	}

	contentDisposition := string(req.Request().Header.Peek("Content-Disposition"))
	if contentDisposition == "" {
		return "", "", nil, aerror.NewCritical(ctx, entity.MissContentDispositionErrorID, nil, "в запросе нет заголовка Content-Disposition")
	}

	_, params, err := mime.ParseMediaType(contentDisposition)
	if err != nil {
		return "", "", nil, aerror.NewCritical(ctx, entity.WrongContentDispositionToErrorID, err, "ошибка при парсинге хедера Content-Disposition (%s): %s", contentDisposition, err.Error())
	}

	tpName, ok := params["filename"]
	if !ok {
		return "", "", nil, aerror.NewCritical(ctx, entity.MissFilenameToErrorID, nil, "в хедере Content-Disposition (%s) нет filename", contentDisposition)
	}

	ok = entity.RxTpName.MatchString(tpName)
	if !ok {
		return "", "", nil, aerror.NewCritical(ctx, entity.WrongFileNameErrorID, nil, "filename %s не соответствует формату", tpName)
	}

	tpBytes := req.Request().Body()
	if len(tpBytes) == 0 {
		return "", "", nil, aerror.NewCritical(ctx, entity.EmptyBodyErrorID, nil, "у запроса пустое тело")
	}

	return tpName, receiptURL, tpBytes, nil
}
