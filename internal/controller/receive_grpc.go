package controller

import (
	"github.com/gofiber/fiber/v2"

	"github.com/efremovich/data-receiver/internal/entity"
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
