package order_repository

import (
	"context"
	"encoding/json"
	"fmt"

	errors_utils "github.com/andreparelho/order-api/pkg/errors"
	"github.com/andreparelho/order-api/pkg/sqs"
	sqs_types "github.com/andreparelho/order-api/pkg/sqs/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type OrderEventRepository interface {
	SendOrderEventMessage(ctx context.Context, queueURL string, event sqs_types.EventOrderCreatedMessage) error
	GetPaymentsOrdersMessage(ctx context.Context, queueURL string) (sqs_types.EventPaymentMessage, types.Message, error)
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

func (o *orderEvent) SendOrderEventMessage(ctx context.Context, queueURL string, event sqs_types.EventOrderCreatedMessage) error {
	orderEventMarsh, err := json.Marshal(&event)
	if err != nil {
		fmt.Printf("\nERROR: erro ao realizar o marshal do event. Erro: %v", err)
		return errors_utils.ErrMarshalEvent
	}

	err = o.sqs.SendMessage(ctx, queueURL, string(orderEventMarsh))
	if err != nil {
		fmt.Printf("\nERROR: erro ao enviar mensagem para fila. Erro: %v", err)
		return errors_utils.ErrSendMessageQueue
	}

	return nil
}

func (o *orderEvent) GetPaymentsOrdersMessage(ctx context.Context, queueURL string) (sqs_types.EventPaymentMessage, types.Message, error) {
	messages, err := o.sqs.ReceiveMessage(ctx, queueURL)
	if err != nil {
		fmt.Printf("\nERROR: erro ao buscar mensagens da fila. Erro: %v", err)
		return sqs_types.EventPaymentMessage{}, types.Message{}, err
	}

	var orderMessages sqs_types.EventPaymentMessage
	var message types.Message
	for _, m := range messages.Messages {
		if err := json.Unmarshal([]byte(*m.Body), &orderMessages); err != nil {
			fmt.Printf("\nERROR: erro ao realizar o unmarshal do event. Erro: %v", err)
			return sqs_types.EventPaymentMessage{}, types.Message{}, err
		}
		message = m
	}

	return orderMessages, message, nil
}

func (o *orderEvent) FinishPaymentOrderEventMessage(ctx context.Context, queueURL string, message *string) error {
	err := o.sqs.DeleteMessage(ctx, queueURL, message)
	if err != nil {
		fmt.Printf("\nERROR: erro ao realizar remocao do evento. Erro: %v", err)
		return err
	}

	return nil
}
