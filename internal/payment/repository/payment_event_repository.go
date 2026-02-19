package payment_event_repository

import (
	"context"
	"encoding/json"
	"fmt"

	errors_utils "github.com/andreparelho/order-api/pkg/errors"
	"github.com/andreparelho/order-api/pkg/sqs"
	sqs_types "github.com/andreparelho/order-api/pkg/sqs/types"
)

type PaymentEventRepostory interface {
	GetOrderPayment(ctx context.Context, queueURL string) (EventOrder, bool, error)
	FinishPaymentProccess(ctx context.Context, queueURL string, message *string) error
	SendPaymentEvent(ctx context.Context, queueURL string, event sqs_types.EventPaymentMessage) error
}

type payment struct {
	sqs sqs.SQSClient
}

func NewPaymentEventRepository(sqs sqs.SQSClient) PaymentEventRepostory {
	return &payment{
		sqs: sqs,
	}
}

type EventOrder struct {
	EventOrderCreatedMessage sqs_types.EventOrderCreatedMessage
	ReceiptHandle            *string
}

func (p *payment) GetOrderPayment(ctx context.Context, queueURL string) (EventOrder, bool, error) {
	messages, err := p.sqs.ReceiveMessage(ctx, queueURL)
	if err != nil {
		fmt.Printf("\n[ERROR]: erro ao buscar mensagens da fila. Erro: %v", err)
		return EventOrder{}, true, err
	}

	if len(messages.Messages) > 0 {
		message := messages.Messages[0]
		var orderMessages sqs_types.EventOrderCreatedMessage
		if err := json.Unmarshal([]byte(*message.Body), &orderMessages); err != nil {
			fmt.Printf("\n[ERROR]: erro ao realizar o unmarshal do event. Erro: %v", err)
			return EventOrder{}, true, err
		}

		eventOrder := EventOrder{
			EventOrderCreatedMessage: orderMessages,
			ReceiptHandle:            message.ReceiptHandle,
		}

		return eventOrder, true, nil
	}

	return EventOrder{}, false, nil
}

func (p *payment) FinishPaymentProccess(ctx context.Context, queueURL string, message *string) error {
	err := p.sqs.DeleteMessage(ctx, queueURL, message)
	if err != nil {
		fmt.Printf("\n[ERROR]: erro eo remover mensagem da fila. Erro: %v", err)
		return err
	}

	return nil
}

func (p *payment) SendPaymentEvent(ctx context.Context, queueURL string, event sqs_types.EventPaymentMessage) error {
	paymentEventMarsh, err := json.Marshal(&event)
	if err != nil {
		fmt.Printf("\n[ERROR]: erro ao realizar o marshal do event. Erro: %v", err)
		return errors_utils.ErrMarshalEvent
	}

	err = p.sqs.SendMessage(ctx, queueURL, string(paymentEventMarsh))
	if err != nil {
		fmt.Printf("\n[ERROR]: erro ao enviar mensagem para fila. Erro: %v", err)
		return errors_utils.ErrSendMessageQueue
	}

	return nil
}
