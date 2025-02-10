package controller

import (
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"

	"github.com/efremovich/data-receiver/internal/entity"
	"github.com/efremovich/data-receiver/pkg/logger"
)

func (gw *grpcGatewayServerImpl) CardReceiveV1Handler(req *fiber.Ctx) error {
	desc := entity.PackageDescription{}

	err := gw.core.ReceiveCards(req.Context(), desc)
	if err != nil {
		return err
	}

	return nil
}

func (gw *grpcGatewayServerImpl) WarehouseReceiveV1Handler(req *fiber.Ctx) error {
	desc := entity.PackageDescription{}

	err := gw.core.ReceiveWarehouses(req.Context(), desc)
	if err != nil {
		return err
	}

	return nil
}

func (gw *grpcGatewayServerImpl) StockReceiverV1Handler(req *fiber.Ctx) error {
	err := gw.core.ReceiveStocks(req.Context(), entity.PackageDescription{})
	if err != nil {
		return err
	}

	return nil
}

func (gw *grpcGatewayServerImpl) OfferFeedV1Handler(req *fiber.Ctx) error {
	filePath := "/tmp/offer_feed.xml"

	logger.GetLoggerFromContext(req.Context()).Infof("начинаем чтение с диска")

	data, err := os.ReadFile(filePath)
	if err != nil {
		req.Response().SetStatusCode(http.StatusInternalServerError)
		req.Response().AppendBodyString("файл с фидом каталога не обнаружен повторите попытку позже")
		return err
	}

	req.Set("Content-Type", "application/xml")
	req.Response().SetStatusCode(http.StatusOK)
	req.Response().AppendBody(data)
	return nil
}
