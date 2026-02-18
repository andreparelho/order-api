package payment_consumer

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	payment_behavior "github.com/andreparelho/order-api/internal/payment/behavior"
	payment_event_repository "github.com/andreparelho/order-api/internal/payment/repository"
	"github.com/andreparelho/order-api/pkg/config"
	sqs_types "github.com/andreparelho/order-api/pkg/sqs/types"
	"github.com/google/uuid"
)

type PaymenytConsumer interface {
	StartWorker(ctx context.Context)
	GetOrdersMessages(ctx context.Context) error
}

type payment struct {
	cfg             config.Configuration
	eventRepository payment_event_repository.PaymentEventRepostory
}

func NewPaymentConsumer(cfg config.Configuration, eventRepository payment_event_repository.PaymentEventRepostory) PaymenytConsumer {
	return &payment{
		cfg:             cfg,
		eventRepository: eventRepository,
	}
}

type PaymentChannel struct {
	ErrorMessage error
}

func (p *payment) StartWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println("[INFO]: worker encerrado")
			return
		default:
			if err := p.GetOrdersMessages(ctx); err != nil {
				fmt.Printf("\n[ERROR]: worker com erro. Erro: %v", err)
			}
		}
	}
}

func (p *payment) GetOrdersMessages(ctx context.Context) error {
	orderEventMessage, message, err := p.eventRepository.GetOrdersPayments(ctx, p.cfg.SQS.OrdersQueue)
	if err != nil {
		fmt.Printf("\n[ERROR]: erro ao buscar as mensagens da fila. Erro: %v", err)
		return err
	}

	return orderPaymentProccess(ctx, orderEventMessage, message.ReceiptHandle, p.cfg.SQS.OrdersQueue, p.cfg.SQS.PaymentsQueue, p.eventRepository)
}

func orderPaymentProccess(ctx context.Context, orderMessage sqs_types.EventOrderCreatedMessage, message *string, orderQueue, paymentQueue string, eventRepository payment_event_repository.PaymentEventRepostory) error {
	err := eventRepository.FinishPaymentProccess(ctx, orderQueue, message)
	if err != nil {
		return err
	}

	eventID, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	paymentEvent := sqs_types.EventPaymentMessage{
		EventId:     fmt.Sprintf("event:payment:{%v}", eventID.String()),
		OrderID:     orderMessage.Data.OrderID,
		EventType:   "payment_completed",
		OccuredTime: time.Now(),
		OrderStatus: string(getPaymentStatus()),
		RedisKey:    orderMessage.Data.RedisKey,
	}

	err = eventRepository.SendPaymentEvent(ctx, paymentQueue, paymentEvent)
	if err != nil {
		return err
	}

	fmt.Printf("\n[INFO]: enviando mensagem para a fila (%s). Evento: payment_completed. Mensagem: %v", paymentQueue, paymentEvent)
	return nil
}

func getPaymentStatus() payment_behavior.PaymentStatus {
	if rand.IntN(100) == 0 {
		return payment_behavior.PaymentStatusCompleted
	}

	return payment_behavior.PaymentStatusFailed
}
