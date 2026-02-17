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
	GetOrdersPayments(ctx context.Context, queueURL string) (sqs_types.EventOrderCreatedMessage, error)
	FinishPaymentProccess(ctx context.Context, queueURL string, message sqs_types.EventOrderCreatedMessage) error
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

func (p *payment) GetOrdersPayments(ctx context.Context, queueURL string) (sqs_types.EventOrderCreatedMessage, error) {
	messages, err := p.sqs.ReceiveMessage(ctx, queueURL)
	if err != nil {
		return sqs_types.EventOrderCreatedMessage{}, err
	}

	var orderMessages sqs_types.EventOrderCreatedMessage
	for _, message := range messages.Messages {
		if err := json.Unmarshal([]byte(*message.Body), &orderMessages); err != nil {
			return sqs_types.EventOrderCreatedMessage{}, err
		}
	}

	return orderMessages, nil
}

func (p *payment) FinishPaymentProccess(ctx context.Context, queueURL string, event sqs_types.EventOrderCreatedMessage) error {
	messageMarsh, err := json.Marshal(event)
	if err != nil {
		return err
	}
	message := string(messageMarsh)

	err = p.sqs.DeleteMessage(ctx, queueURL, &message)
	if err != nil {
		return err
	}

	return nil
}

func (p *payment) SendPaymentEvent(ctx context.Context, queueURL string, event sqs_types.EventPaymentMessage) error {
	paymentEventMarsh, err := json.Marshal(&event)
	if err != nil {
		fmt.Printf("ERROR: erro ao realizar o marshal do event, erro: %v", err)
		return errors_utils.ErrMarshalEvent
	}

	err = p.sqs.SendMessage(ctx, queueURL, string(paymentEventMarsh))
	if err != nil {
		fmt.Printf("ERROR: erro ao enviar mensagem para fila, erro: %v", err)
		return errors_utils.ErrSendMessageQueue
	}

	return nil
}
