package order_handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	order_service "github.com/andreparelho/order-api/internal/order/service"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/google/uuid"
)

func OrderHandler(orderService order_service.OrderService) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		xRequestId := requestid.FromContext(ctx)
		if xRequestId == "" {
			id, err := uuid.NewRandom()
			if err != nil {
				fmt.Printf("ERROR: erro ao criar um novo x-request-id. Erro: %v", err)
				return ctx.SendStatus(http.StatusInternalServerError)
			}
			xRequestId = id.String()
		}

		var orderRequest order_service.CreateOrderRequest
		if err := json.Unmarshal(ctx.Body(), &orderRequest); err != nil {
			fmt.Printf("ERROR: erro ao realizar o unmarshal do request, erro: %v", err)
			return ctx.SendStatus(http.StatusBadRequest)
		}

		err := orderService.CreateOrderService(ctx.Context(), orderRequest, xRequestId)
		if err != nil {
			fmt.Printf("ERROR: erro interno, erro: %v", err)
			return ctx.SendStatus(http.StatusInternalServerError)
		}

		return ctx.SendStatus(http.StatusCreated)
	}
}
