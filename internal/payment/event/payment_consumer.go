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
	StartConsumer(ctx context.Context)
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

func (p *payment) StartConsumer(ctx context.Context) {
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
	orderEventMessage, haveMessage, err := p.eventRepository.GetOrderPayment(ctx, p.cfg.SQS.OrdersQueue)
	if err != nil {
		fmt.Printf("\n[ERROR]: erro ao buscar as mensagens da fila. Erro: %v", err)
		return err
	}

	if !haveMessage {
		fmt.Print("\n[INFO]: nenhuma mensagem na fila")
		return nil
	}

	return orderPaymentProccess(ctx, orderEventMessage, p.cfg.SQS.OrdersQueue, p.cfg.SQS.PaymentsQueue, p.eventRepository)
}

func orderPaymentProccess(ctx context.Context, orderMessage payment_event_repository.EventOrder, orderQueue, paymentQueue string, eventRepository payment_event_repository.PaymentEventRepostory) error {
	receiptHandle := orderMessage.ReceiptHandle

	err := eventRepository.FinishPaymentProccess(ctx, orderQueue, receiptHandle)
	if err != nil {
		return err
	}
	fmt.Printf("\n[INFO]: deletando mensagem da fila (%s). Mensagem: %v. Source: payment_service", orderQueue, orderMessage)

	eventID, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	paymentEvent := sqs_types.EventPaymentMessage{
		EventId:     fmt.Sprintf("event:payment:{%s}", eventID.String()),
		OrderID:     orderMessage.EventOrderCreatedMessage.Data.OrderID,
		EventType:   "payment_completed",
		Source:      "payment_service",
		OccuredTime: time.Now(),
		OrderStatus: string(getPaymentStatus()),
		CacheKey:    orderMessage.EventOrderCreatedMessage.Data.CacheKey,
	}

	err = eventRepository.SendPaymentEvent(ctx, paymentQueue, paymentEvent)
	if err != nil {
		return err
	}

	fmt.Printf("\n[INFO]: enviando mensagem para a fila (%s). Evento: payment_completed. Mensagem: %v. Source: payment_service", paymentQueue, paymentEvent)
	return nil
}

func getPaymentStatus() payment_behavior.PaymentStatus {
	if rand.IntN(10)%2 == 0 {
		return payment_behavior.PaymentStatusCompleted
	}

	return payment_behavior.PaymentStatusFailed
}
