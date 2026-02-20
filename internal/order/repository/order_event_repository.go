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
	GetPaymentOrderMessage(ctx context.Context, queueURL string) (EventPayment, bool, error)
	FinishPaymentOrderEventMessage(ctx context.Context, queueURL string, message *string) error
}

type orderEvent struct {
	sqs sqs.SQSClient
}

func NewOrderEventRepository(sqs sqs.SQSClient) OrderEventRepository {
	return &orderEvent{
		sqs: sqs,
	}
}

type EventPayment struct {
	EventPaymentMessage sqs_types.EventPaymentMessage
	ReceiptHandle       *string
}

func (o *orderEvent) SendOrderEventMessage(ctx context.Context, queueURL string, event sqs_types.EventOrderCreatedMessage) error {
	orderEventMarsh, err := json.Marshal(&event)
	if err != nil {
		fmt.Printf("\n[ERROR]: erro ao realizar o marshal do event. Erro: %v", err)
		return errors_utils.ErrMarshalEvent
	}

	if err := o.sqs.SendMessage(ctx, queueURL, string(orderEventMarsh)); err != nil {
		fmt.Printf("\n[ERROR]: erro ao enviar mensagem para fila. Erro: %v", err)
		return errors_utils.ErrSendMessageQueue
	}

	return nil
}

func (o *orderEvent) GetPaymentOrderMessage(ctx context.Context, queueURL string) (EventPayment, bool, error) {
	messages, err := o.sqs.ReceiveMessage(ctx, queueURL)
	if err != nil {
		fmt.Printf("\n[ERROR]: erro ao buscar mensagens da fila. Erro: %v", err)
		return EventPayment{}, true, err
	}

	if len(messages.Messages) > 0 {
		message := messages.Messages[0]
		var paymentMessages sqs_types.EventPaymentMessage
		if err := json.Unmarshal([]byte(*message.Body), &paymentMessages); err != nil {
			fmt.Printf("\n[ERROR]: erro ao realizar o unmarshal do event. Erro: %v", err)
			return EventPayment{}, true, err
		}

		eventPayment := EventPayment{
			EventPaymentMessage: paymentMessages,
			ReceiptHandle:       message.ReceiptHandle,
		}

		return eventPayment, true, nil
	}

	return EventPayment{}, false, nil
}

func (o *orderEvent) FinishPaymentOrderEventMessage(ctx context.Context, queueURL string, message *string) error {
	if err := o.sqs.DeleteMessage(ctx, queueURL, message); err != nil {
		fmt.Printf("\n[ERROR]: erro ao realizar remocao do evento. Erro: %v", err)
		return err
	}

	return nil
}
