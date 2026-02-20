package payment_repository

import (
	"context"
	"strconv"
	"time"

	"github.com/andreparelho/order-api/pkg/dynamo"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
)

type PaymentRepository interface {
	SaveOrderPayment(ctx context.Context, orderPayment OrderPayment) error
}

type paymentDatabase struct {
	dynamoClient dynamo.DynamoClient
	tableName    string
}

type OrderPayment struct {
	OrderID   uuid.UUID
	PaymentID uuid.UUID
	EventID   string
	Status    string
	Amount    float64
	Currency  string
	CreatedAt time.Time
}

func NewPaymentRepository(dynamoClient dynamo.DynamoClient, tableName string) PaymentRepository {
	return &paymentDatabase{
		dynamoClient: dynamoClient,
		tableName:    tableName,
	}
}

func (p *paymentDatabase) SaveOrderPayment(ctx context.Context, orderPayment OrderPayment) error {
	item := map[string]types.AttributeValue{
		"order_id":   &types.AttributeValueMemberS{Value: orderPayment.OrderID.String()},
		"payment_id": &types.AttributeValueMemberS{Value: orderPayment.PaymentID.String()},
		"event_id":   &types.AttributeValueMemberS{Value: orderPayment.EventID},
		"status":     &types.AttributeValueMemberS{Value: orderPayment.Status},
		"amount":     &types.AttributeValueMemberN{Value: strconv.FormatFloat(orderPayment.Amount, 'f', -1, 64)},
		"currency":   &types.AttributeValueMemberS{Value: orderPayment.Currency},
		"created_at": &types.AttributeValueMemberS{Value: orderPayment.CreatedAt.Format(time.RFC3339)},
	}

	if err := p.dynamoClient.PutItem(ctx, item, p.tableName); err != nil {
		return err
	}

	return nil
}
