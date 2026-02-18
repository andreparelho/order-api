package payment_event_repository

import (
	"context"
	"encoding/json"
	"fmt"

	errors_utils "github.com/andreparelho/order-api/pkg/errors"
	"github.com/andreparelho/order-api/pkg/sqs"
	sqs_types "github.com/andreparelho/order-api/pkg/sqs/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type PaymentEventRepostory interface {
	GetOrdersPayments(ctx context.Context, queueURL string) (sqs_types.EventOrderCreatedMessage, types.Message, error)
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

func (p *payment) GetOrdersPayments(ctx context.Context, queueURL string) (sqs_types.EventOrderCreatedMessage, types.Message, error) {
	messages, err := p.sqs.ReceiveMessage(ctx, queueURL)
	if err != nil {
		fmt.Printf("\nERROR: erro ao buscar mensagens da fila. Erro: %v", err)
		return sqs_types.EventOrderCreatedMessage{}, types.Message{}, err
	}

	var orderMessages sqs_types.EventOrderCreatedMessage
	var message types.Message
	for _, m := range messages.Messages {
		if err := json.Unmarshal([]byte(*m.Body), &orderMessages); err != nil {
			fmt.Printf("\nERROR: erro ao realizar o unmarshal do event. Erro: %v", err)
			return sqs_types.EventOrderCreatedMessage{}, types.Message{}, err
		}
		message = m
	}

	return orderMessages, message, nil
}

func (p *payment) FinishPaymentProccess(ctx context.Context, queueURL string, message *string) error {
	err := p.sqs.DeleteMessage(ctx, queueURL, message)
	if err != nil {
		fmt.Printf("\nERROR: erro eo remover mensagem da fila. Erro: %v", err)
		return err
	}

	return nil
}

func (p *payment) SendPaymentEvent(ctx context.Context, queueURL string, event sqs_types.EventPaymentMessage) error {
	paymentEventMarsh, err := json.Marshal(&event)
	if err != nil {
		fmt.Printf("\nERROR: erro ao realizar o marshal do event. Erro: %v", err)
		return errors_utils.ErrMarshalEvent
	}

	err = p.sqs.SendMessage(ctx, queueURL, string(paymentEventMarsh))
	if err != nil {
		fmt.Printf("\nERROR: erro ao enviar mensagem para fila. Erro: %v", err)
		return errors_utils.ErrSendMessageQueue
	}

	return nil
}
