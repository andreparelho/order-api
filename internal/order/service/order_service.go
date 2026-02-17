package order_service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	order_behavior "github.com/andreparelho/order-api/internal/order/behavior"
	order_repository "github.com/andreparelho/order-api/internal/order/repository"
	"github.com/andreparelho/order-api/pkg/config"
	"github.com/andreparelho/order-api/pkg/sqs"
	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrderService(ctx context.Context, orderRequest CreateOrderRequest, xRequestId string) error
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

var (
	ErrGenerateUUID     = fmt.Errorf("erro criar um uuid")
	ErrDatabaseInsert   = fmt.Errorf("erro ao inserir o dado na base")
	ErrSendMessageQueue = fmt.Errorf("erro ao enviar mensagem para fila")
	ErrMarshalEvent     = fmt.Errorf("erro ao realizar o marshal do evento sqs")
)

type CreateOrderRequest struct {
	CustomerID  uuid.UUID `json:"customerID"`
	TotalAmount float64   `json:"totalAmount"`
	Currency    string    `json:"currency"`
}

type EventOrderCreatedMessage struct {
	EventId     string         `json:"eventID"`
	EventType   string         `json:"eventType"`
	OccuredTime time.Time      `json:"occuredTime"`
	Data        OrderEventData `json:"data"`
}

type OrderEventData struct {
	OrderID     uuid.UUID `json:"orderID"`
	CustomerID  uuid.UUID `json:"customerID"`
	TotalAmount float64   `json:"totalAmount"`
	Currency    string    `json:"currency"`
}

func (o *order) CreateOrderService(ctx context.Context, orderRequest CreateOrderRequest, xRequestId string) error {
	orderID, err := uuid.NewRandom()
	if err != nil {
		fmt.Printf("ERROR: erro criar um uuid, erro: %v", err)
		return ErrGenerateUUID
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

	isRedisOk, err := o.repository.InsertOrder(ctx, order, xRequestId)
	if err != nil {
		fmt.Printf("ERROR: erro ao inserir o dado na base, erro: %v", err)
		return ErrDatabaseInsert
	} else if isRedisOk {
		fmt.Printf("INFO: ordem encontrada no redis. Encerrando fluxo %v", order)
		return nil
	}

	eventID, err := uuid.NewRandom()
	if err != nil {
		fmt.Printf("ERROR: erro criar um uuid, erro: %v", err)
		return ErrGenerateUUID
	}

	orderEvent := EventOrderCreatedMessage{
		EventId:     fmt.Sprintf("event:created:{%s}", eventID.String()),
		EventType:   "orderCreated",
		OccuredTime: time.Now(),
		Data: OrderEventData{
			OrderID:     orderID,
			CustomerID:  order.CustomerID,
			TotalAmount: order.TotalAmount,
			Currency:    order.Currency,
		},
	}

	orderEventMarsh, err := json.Marshal(&orderEvent)
	if err != nil {
		fmt.Printf("ERROR: erro ao realizar o marshal do event, erro: %v", err)
		return ErrMarshalEvent
	}

	err = o.sqs.SendMessage(ctx, o.cfg.SQS.QueueName, string(orderEventMarsh))
	if err != nil {
		fmt.Printf("ERROR: erro ao enviar mensagem para fila, erro: %v", err)
		return ErrSendMessageQueue
	}

	return nil
}
