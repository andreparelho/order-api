package order_service

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	order_behavior "github.com/andreparelho/order-api/internal/order/behavior"
	order_repository "github.com/andreparelho/order-api/internal/order/repository"
	"github.com/andreparelho/order-api/pkg/config"
	errors_utils "github.com/andreparelho/order-api/pkg/errors"
	"github.com/andreparelho/order-api/pkg/sqs"
	sqs_types "github.com/andreparelho/order-api/pkg/sqs/types"
	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrderService(ctx context.Context, orderRequest CreateOrderRequest, xRequestId string) error
}

type order struct {
	repository      order_repository.OrderRepository
	eventRepository order_repository.OrderEventRepository
	sqs             sqs.SQSClient
	cfg             config.Configuration
}

func NewOrderService(orderRepository order_repository.OrderRepository, eventRepository order_repository.OrderEventRepository, cfg config.Configuration) OrderService {
	return &order{
		repository:      orderRepository,
		eventRepository: eventRepository,
		cfg:             cfg,
	}
}

type CreateOrderRequest struct {
	CustomerID  uuid.UUID `json:"customerID"`
	TotalAmount float64   `json:"totalAmount"`
	Currency    string    `json:"currency"`
}

func (o *order) CreateOrderService(ctx context.Context, orderRequest CreateOrderRequest, xRequestId string) error {
	orderID, err := uuid.NewRandom()
	if err != nil {
		fmt.Printf("\n[ERROR]: erro criar um uuid, erro: %v", err)
		return errors_utils.ErrGenerateUUID
	}

	order := order_repository.Order{
		OrderID:     orderID,
		CustomerID:  orderRequest.CustomerID,
		Status:      string(order_behavior.OrderStatusCreated),
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

	redisKey := fmt.Sprintf("order:req:id_%s", xRequestId)

	isRedisOk, err := o.repository.InsertOrder(ctx, order, redisKey)
	if err != nil {
		fmt.Printf("\n[ERROR]: erro ao inserir o dado na base, erro: %v", err)
		return errors_utils.ErrDatabaseInsert
	} else if isRedisOk {
		fmt.Printf("\n[INFO]: ordem encontrada no redis. Encerrando fluxo %v", order)
		return nil
	}

	eventID, err := uuid.NewRandom()
	if err != nil {
		fmt.Printf("\n[ERROR]: erro criar um uuid, erro: %v", err)
		return errors_utils.ErrGenerateUUID
	}

	orderEvent := sqs_types.EventOrderCreatedMessage{
		EventId:     fmt.Sprintf("event:order_created:{%s}", eventID.String()),
		EventType:   "order_created",
		OccuredTime: time.Now(),
		Data: sqs_types.OrderEventData{
			OrderID:     orderID,
			CustomerID:  order.CustomerID,
			RedisKey:    redisKey,
			TotalAmount: order.TotalAmount,
			Currency:    order.Currency,
		},
	}

	err = o.eventRepository.SendOrderEventMessage(ctx, o.cfg.SQS.OrdersQueue, orderEvent)
	if err != nil {
		fmt.Printf("\n[ERROR]: erro ao enviar mensagem para fila, erro: %v", err)
		return errors_utils.ErrSendMessageQueue
	}
	fmt.Printf("\n[INFO]: enviando mensagem para a fila (%s). Evento: order_created. Mensagem: %v", o.cfg.SQS.OrdersQueue, orderEvent)

	return nil
}
