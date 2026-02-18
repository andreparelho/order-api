package order_consumer

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	order_repository "github.com/andreparelho/order-api/internal/order/repository"
	"github.com/andreparelho/order-api/pkg/config"
	sqs_types "github.com/andreparelho/order-api/pkg/sqs/types"
)

type OrderConsumer interface {
	StartWorker(ctx context.Context)
	GetPaymentsMessages(ctx context.Context) error
}

type order struct {
	cfg             config.Configuration
	eventRepository order_repository.OrderEventRepository
	orderRepository order_repository.OrderRepository
}

func NewOrderConsumer(cfg config.Configuration, eventRepository order_repository.OrderEventRepository, orderRepository order_repository.OrderRepository) OrderConsumer {
	return &order{
		cfg:             cfg,
		eventRepository: eventRepository,
		orderRepository: orderRepository,
	}
}

type OrderChannel struct {
	ErrorMessage error
}

func (o *order) StartWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("[INFO]: worker encerrado")
			return
		default:
			if err := o.GetPaymentsMessages(ctx); err != nil {
				fmt.Printf("\n[ERROR]: worker com erro. Erro: %v", err)
			}
		}
	}
}

func (o *order) GetPaymentsMessages(ctx context.Context) error {
	paymentOrderEventMessage, message, err := o.eventRepository.GetPaymentsOrdersMessage(ctx, o.cfg.SQS.PaymentsQueue)
	if err != nil {
		fmt.Printf("\n[ERROR]: erro ao buscar as mensagens da fila. Erro: %v", err)
		return err
	}

	return finishOrderProccess(ctx, paymentOrderEventMessage, message.ReceiptHandle, o.cfg.SQS.PaymentsQueue, o.eventRepository, o.orderRepository)
}

func finishOrderProccess(ctx context.Context, paymentMessage sqs_types.EventPaymentMessage, message *string, paymentQueue string, eventRepository order_repository.OrderEventRepository, orderRepository order_repository.OrderRepository) error {
	order := order_repository.Order{
		OrderID: paymentMessage.OrderID,
		Status:  paymentMessage.OrderStatus,
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}

	err := orderRepository.UpdateOrder(ctx, order, paymentMessage.RedisKey)
	if err != nil {
		return err
	}

	err = eventRepository.FinishPaymentOrderEventMessage(ctx, paymentQueue, message)
	if err != nil {
		return err
	}

	fmt.Printf("\n[INFO]: deletando mensagem da fila (%s). Mensagem: %v", paymentQueue, message)
	return nil
}
