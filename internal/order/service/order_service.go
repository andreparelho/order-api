package order_service

import (
	"encoding/json"
	"net/http"

	order_repository "github.com/andreparelho/order-api/internal/order/repository"
	"github.com/andreparelho/order-api/pkg/redis"
	"github.com/gofiber/fiber/v3"
)

type OrderService interface {
	CreateOrder() fiber.Handler
}

type order struct {
	Redis      redis.RedisClient
	Repository order_repository.OrderRepository
}

func NewOrderService(redis redis.RedisClient, orderRepository order_repository.OrderRepository) OrderService {
	return &order{
		Redis:      redis,
		Repository: orderRepository,
	}
}

type CreateOrderRequest struct {
	ProductName string  `json:"productName"`
	Value       float32 `json:"value"`
}

func (o *order) CreateOrder() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		var orderRequest CreateOrderRequest
		if err := json.Unmarshal(ctx.Body(), &orderRequest); err != nil {
			return ctx.SendStatus(http.StatusBadRequest)
		}

		return ctx.SendStatus(http.StatusCreated)
	}
}
