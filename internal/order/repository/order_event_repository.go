package order_repository

import (
	"context"
	"encoding/json"
	"fmt"

	errors_utils "github.com/andreparelho/order-api/pkg/errors"
	"github.com/andreparelho/order-api/pkg/sqs"
	sqs_types "github.com/andreparelho/order-api/pkg/sqs/types"
)

type OrderEventRepository interface {
	SendOrderEventMessage(ctx context.Context, queueURL string, event sqs_types.EventOrderCreatedMessage) error
	GetPaymentsOrdersMessage(ctx context.Context, queueURL string) (sqs_types.EventPaymentMessage, error)
	FinishPaymentOrderEventMessage(ctx context.Context, queueURL string, event sqs_types.EventPaymentMessage) error
}

type orderEvent struct {
	sqs sqs.SQSClient
}

func NewOrderEventRepository(sqs sqs.SQSClient) OrderEventRepository {
	return &orderEvent{
		sqs: sqs,
	}
}

func (o *orderEvent) SendOrderEventMessage(ctx context.Context, queueURL string, event sqs_types.EventOrderCreatedMessage) error {
	orderEventMarsh, err := json.Marshal(&event)
	if err != nil {
		fmt.Printf("ERROR: erro ao realizar o marshal do event, erro: %v", err)
		return errors_utils.ErrMarshalEvent
	}

	err = o.sqs.SendMessage(ctx, queueURL, string(orderEventMarsh))
	if err != nil {
		fmt.Printf("ERROR: erro ao enviar mensagem para fila, erro: %v", err)
		return errors_utils.ErrSendMessageQueue
	}

	return nil
}

func (o *orderEvent) GetPaymentsOrdersMessage(ctx context.Context, queueURL string) (sqs_types.EventPaymentMessage, error) {
	messages, err := o.sqs.ReceiveMessage(ctx, queueURL)
	if err != nil {
		return sqs_types.EventPaymentMessage{}, err
	}

	var orderMessages sqs_types.EventPaymentMessage
	for _, message := range messages.Messages {
		if err := json.Unmarshal([]byte(*message.Body), &orderMessages); err != nil {
			return sqs_types.EventPaymentMessage{}, err
		}
	}

	return orderMessages, nil
}

func (o *orderEvent) FinishPaymentOrderEventMessage(ctx context.Context, queueURL string, event sqs_types.EventPaymentMessage) error {
	messageMarsh, err := json.Marshal(event)
	if err != nil {
		return err
	}
	message := string(messageMarsh)

	err = o.sqs.DeleteMessage(ctx, queueURL, &message)
	if err != nil {
		return err
	}

	return nil
}
