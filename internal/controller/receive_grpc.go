package controller

import (
	"net/http"
	"strconv"

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

func (gw *grpcGatewayServerImpl) OfferFeedV1Handler(req *fiber.Ctx) error {
	data, err := gw.core.OfferFeed(req.Context())

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

func (gw *grpcGatewayServerImpl) StockFeedV1Handler(req *fiber.Ctx) error {
	data, err := gw.core.StockFeed(req.Context())
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

func (gw *grpcGatewayServerImpl) VKCardsFeedV1Handler(req *fiber.Ctx) error {
	// Извлечь параметры из запроса
	queryLimit := req.Query("limit", "0")
	limit, err := strconv.Atoi(queryLimit)
	if err != nil {
		limit = 0
	}

	cursor := req.Query("cursor", "")
	filter := req.Query("filter", "")

	params := entity.VkCardsFeedParams{
		Limit:  limit,
		Cursor: cursor,
		Filter: filter,
	}
	data, err := gw.core.VkCardsFeed(req.Context(), params)
	if err != nil {
		req.Response().SetStatusCode(http.StatusInternalServerError)
		req.Response().AppendBodyString("файл с фидом каталога не обнаружен повторите попытку позже")
		return err
	}
	req.Set("Content-Type", "application/json")
	req.Response().SetStatusCode(http.StatusOK)
	req.Response().AppendBody(data)
	return nil
}
