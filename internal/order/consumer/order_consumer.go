package order_consumer

import (
	"context"

	order_repository "github.com/andreparelho/order-api/internal/order/repository"
	"github.com/andreparelho/order-api/pkg/config"
	sqs_types "github.com/andreparelho/order-api/pkg/sqs/types"
)

type OrderConsumer interface {
	GetPaymentsMessages(ctx context.Context) error
}

type order struct {
	cfg             config.Configuration
	eventRepository order_repository.OrderEventRepository
	orderRepository order_repository.OrderRepository
}

type OrderChannel struct {
	ErrorMessage error
}

func (o *order) GetPaymentsMessages(ctx context.Context) error {
	for {
		paymentOrderEventMessage, err := o.eventRepository.GetPaymentsOrdersMessage(ctx, o.cfg.SQS.PaymentsQueue)
		if err != nil {
			return err
		}

		if paymentOrderEventMessage.EventType != "payment_completed" {
			continue
		}

		mch := make(chan OrderChannel)
		go finishOrderProccess(ctx, paymentOrderEventMessage, o.cfg.SQS.PaymentsQueue, o.eventRepository, o.orderRepository, mch)

		errMch := <-mch
		if errMch.ErrorMessage != nil {
			return err
		}
	}
}

func finishOrderProccess(ctx context.Context, message sqs_types.EventPaymentMessage, queueUrl string, eventRepository order_repository.OrderEventRepository, orderRepository order_repository.OrderRepository, ch chan OrderChannel) {
	defer close(ch)

	err := orderRepository.UpdateOrder(ctx, order_repository.Order{})
	if err != nil {
		ch <- OrderChannel{ErrorMessage: err}
		return
	}

	err = eventRepository.FinishPaymentOrderEventMessage(ctx, queueUrl, message)
	if err != nil {
		ch <- OrderChannel{ErrorMessage: err}
		return
	}
}
