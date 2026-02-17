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

func (p *payment) GetOrdersMessages(ctx context.Context) error {
	for {
		orderEventMessage, err := p.eventRepository.GetOrdersPayments(ctx, p.cfg.SQS.OrdersQueue)
		if err != nil {
			return err
		}

		if orderEventMessage.EventType != "order_created" {
			continue
		}

		mch := make(chan PaymentChannel)
		go orderPaymentProccess(ctx, orderEventMessage, p.cfg.SQS.OrdersQueue, p.cfg.SQS.PaymentsQueue, p.eventRepository, mch)

		errMch := <-mch
		if errMch.ErrorMessage != nil {
			return err
		}
	}
}

func orderPaymentProccess(ctx context.Context, message sqs_types.EventOrderCreatedMessage, orderQueue, paymentQueue string, eventRepository payment_event_repository.PaymentEventRepostory, ch chan PaymentChannel) {
	defer close(ch)

	err := eventRepository.FinishPaymentProccess(ctx, orderQueue, message)
	if err != nil {
		ch <- PaymentChannel{ErrorMessage: err}
		return
	}

	eventID, err := uuid.NewRandom()
	if err != nil {
		ch <- PaymentChannel{ErrorMessage: err}
		return
	}

	paymentEvent := sqs_types.EventPaymentMessage{
		EventId:     fmt.Sprintf("event:payment:{%v}", eventID.String()),
		EventType:   "payment_completed",
		OccuredTime: time.Now(),
		OrderStatus: string(getPaymentStatus()),
	}

	err = eventRepository.SendPaymentEvent(ctx, paymentQueue, paymentEvent)
	if err != nil {
		ch <- PaymentChannel{ErrorMessage: err}
		return
	}
}

func getPaymentStatus() payment_behavior.PaymentStatus {
	if rand.IntN(100) == 0 {
		return payment_behavior.PaymentStatusCompleted
	}

	return payment_behavior.PaymentStatusFailed
}
