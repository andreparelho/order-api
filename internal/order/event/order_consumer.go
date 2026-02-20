package order_consumer

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	order_repository "github.com/andreparelho/order-api/internal/order/repository"
	"github.com/andreparelho/order-api/pkg/config"
)

type OrderConsumer interface {
	StartConsumer(ctx context.Context)
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

func (o *order) StartConsumer(ctx context.Context) {
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
	paymentOrderEventMessage, haveMessage, err := o.eventRepository.GetPaymentOrderMessage(ctx, o.cfg.SQS.PaymentsQueue)
	if err != nil {
		fmt.Printf("\n[ERROR]: erro ao buscar as mensagens da fila. Erro: %v", err)
		return err
	}

	if !haveMessage {
		return nil
	}

	return finishOrderProccess(ctx, paymentOrderEventMessage, o.cfg.SQS.PaymentsQueue, o.eventRepository, o.orderRepository)
}

func finishOrderProccess(ctx context.Context, paymentMessage order_repository.EventPayment, paymentQueue string, eventRepository order_repository.OrderEventRepository, orderRepository order_repository.OrderRepository) error {
	order := order_repository.Order{
		OrderID: paymentMessage.EventPaymentMessage.OrderID,
		Status:  paymentMessage.EventPaymentMessage.OrderStatus,
		UpdatedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}

	time.Sleep(10 * time.Second)

	err := orderRepository.UpdateOrder(ctx, order, paymentMessage.EventPaymentMessage.CacheKey)
	if err != nil {
		return err
	}

	err = eventRepository.FinishPaymentOrderEventMessage(ctx, paymentQueue, paymentMessage.ReceiptHandle)
	if err != nil {
		return err
	}

	fmt.Printf("\n[INFO]: deletando mensagem da fila (%s). Mensagem: %v", paymentQueue, paymentMessage)
	return nil
}
