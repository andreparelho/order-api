package order_service

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	order_behavior "github.com/andreparelho/order-api/internal/order/behavior"
	order_repository "github.com/andreparelho/order-api/internal/order/repository"
	"github.com/andreparelho/order-api/pkg/config"
	"github.com/andreparelho/order-api/pkg/sqs"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder() fiber.Handler
}

type order struct {
	repository order_repository.OrderRepository
	sqs        sqs.SQSClient
	cfg        config.Configuration
}

func NewOrderService(orderRepository order_repository.OrderRepository, sqs sqs.SQSClient, cfg config.Configuration) OrderService {
	return &order{
		repository: orderRepository,
		sqs:        sqs,
		cfg:        cfg,
	}
}

type CreateOrderRequest struct {
	CustomerID  uuid.UUID `json:"customerID"`
	TotalAmount float64   `json:"totalAmount"`
	Currency    string    `json:"currency"`
}

func (o *order) CreateOrder() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		var orderRequest CreateOrderRequest
		if err := json.Unmarshal(ctx.Body(), &orderRequest); err != nil {
			fmt.Printf("erro ao realizar o unmarshal do request, erro: %v", err)
			return ctx.SendStatus(http.StatusBadRequest)
		}

		orderID, err := uuid.NewRandom()
		if err != nil {
			fmt.Printf("erro criar um uuid, erro: %v", err)
			return ctx.SendStatus(http.StatusInternalServerError)
		}

		order := order_repository.Order{
			OrderID:     orderID,
			CustomerID:  orderRequest.CustomerID,
			Status:      order_behavior.OrderStatusCreated,
			TotalAmount: orderRequest.TotalAmount,
			Currency:    orderRequest.Currency,
			CreatedAt: sql.NullTime{
				Time:  time.Now(),
				Valid: true,
			},
			UpdatedAt: sql.NullTime{
				Time:  time.Time{},
				Valid: false,
			},
		}

		err = o.repository.InsertOrder(ctx.Context(), order)
		if err != nil {
			fmt.Printf("erro ao inserir o dado na base, erro: %v", err)
			return ctx.SendStatus(http.StatusInternalServerError)
		}

		err = o.sqs.SendMessage(ctx.Context(), o.cfg.SQS.QueueName, "teste")
		if err != nil {
			fmt.Printf("erro ao enviar mensagem para fila, erro: %v", err)
			return ctx.SendStatus(http.StatusInternalServerError)
		}

		return ctx.SendStatus(http.StatusCreated)
	}
}
