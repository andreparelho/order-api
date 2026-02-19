package sqs_types

import (
	"time"

	"github.com/google/uuid"
)

type EventOrderCreatedMessage struct {
	EventId     string         `json:"eventID"`
	EventType   string         `json:"eventType"`
	Source      string         `json:"source"`
	OccuredTime time.Time      `json:"occuredTime"`
	Data        OrderEventData `json:"data"`
}

type OrderEventData struct {
	OrderID     uuid.UUID `json:"orderID"`
	CustomerID  uuid.UUID `json:"customerID"`
	CacheKey    string    `json:"cacheKey"`
	TotalAmount float64   `json:"totalAmount"`
	Currency    string    `json:"currency"`
}

type EventPaymentMessage struct {
	EventId     string    `json:"eventID"`
	OrderID     uuid.UUID `json:"orderID"`
	EventType   string    `json:"eventType"`
	Source      string    `json:"source"`
	OccuredTime time.Time `json:"occuredTime"`
	OrderStatus string    `json:"status"`
	CacheKey    string    `json:"cacheKey"`
}
