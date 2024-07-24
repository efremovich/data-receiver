package controller

import (
	"github.com/gofiber/fiber/v2"

	"github.com/efremovich/data-receiver/internal/entity"
)

func (gw *grpcGatewayServerImpl) CardReceiveV1Handler(req *fiber.Ctx) error {
	desc := entity.PackageDescription{
		Cursor: 0,
		Limit:  100,
	}

	err := gw.core.ReceiveCards(req.Context(), desc)
	if err != nil {
		return err
	}
	return nil
}
