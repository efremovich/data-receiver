package controller

import "github.com/gofiber/fiber/v2"

func (gw *grpcGatewayServerImpl) CardReceiveV1Handler(req *fiber.Ctx) error {
	err := gw.core.ReceiveCards(req.Context(), 0)
	if err != nil {
		return err
	}
	return nil
}
