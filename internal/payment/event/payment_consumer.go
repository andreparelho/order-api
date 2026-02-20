package payment_consumer

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	payment_behavior "github.com/andreparelho/order-api/internal/payment/behavior"
	payment_repository "github.com/andreparelho/order-api/internal/payment/repository"
	"github.com/andreparelho/order-api/pkg/config"
	sqs_types "github.com/andreparelho/order-api/pkg/sqs/types"
	"github.com/google/uuid"
)

type PaymenytConsumer interface {
	StartConsumer(ctx context.Context)
	GetOrdersMessages(ctx context.Context) error
}

type payment struct {
	cfg               config.Configuration
	eventRepository   payment_repository.PaymentEventRepostory
	paymentRepository payment_repository.PaymentRepository
}

func NewPaymentConsumer(cfg config.Configuration, eventRepository payment_repository.PaymentEventRepostory, paymentRepository payment_repository.PaymentRepository) PaymenytConsumer {
	return &payment{
		cfg:               cfg,
		eventRepository:   eventRepository,
		paymentRepository: paymentRepository,
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
		return nil
	}

	paymentID, err := uuid.NewRandom()
	if err != nil {
		fmt.Printf("\n[ERROR]: erro eo gerar o paymentID. Erro: %v", err)
		return err
	}

	paymentStatus := getPaymentStatus()

	orderPayment := payment_repository.OrderPayment{
		OrderID:   orderEventMessage.EventOrderCreatedMessage.Data.OrderID,
		PaymentID: paymentID,
		EventID:   orderEventMessage.EventOrderCreatedMessage.EventID,
		Status:    string(paymentStatus),
		Amount:    orderEventMessage.EventOrderCreatedMessage.Data.TotalAmount,
		CreatedAt: time.Now(),
	}

	if err := p.paymentRepository.SaveOrderPayment(ctx, orderPayment); err != nil {
		fmt.Printf("\n[ERROR]: erro ao salvar a ordem de pagamento no banco de dados. Erro: %v", err)
		return err
	}

	return orderPaymentProccess(ctx, orderEventMessage, p.cfg.SQS.OrdersQueue, p.cfg.SQS.PaymentsQueue, p.eventRepository, string(paymentStatus))
}

func orderPaymentProccess(ctx context.Context, orderMessage payment_repository.EventOrder, orderQueue, paymentQueue string, eventRepository payment_repository.PaymentEventRepostory, paymentStatus string) error {
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
		EventID:     fmt.Sprintf("event:payment:{%s}", eventID.String()),
		OrderID:     orderMessage.EventOrderCreatedMessage.Data.OrderID,
		EventType:   "payment_completed",
		Source:      "payment_service",
		OccuredTime: time.Now(),
		OrderStatus: paymentStatus,
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
