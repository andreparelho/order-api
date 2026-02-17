package order_repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	order_behavior "github.com/andreparelho/order-api/internal/order/behavior"
	"github.com/andreparelho/order-api/pkg/redis"
	"github.com/google/uuid"
)

type OrderRepository interface {
	InsertOrder(ctx context.Context, order Order) error
}

type order struct {
	database *sql.DB
	redis    redis.RedisClient
}

func NewOrderRepository(db *sql.DB, redis redis.RedisClient) OrderRepository {
	return &order{
		database: db,
		redis:    redis,
	}
}

type Order struct {
	OrderID     uuid.UUID                  `db:"id" json:"orderID"`
	CustomerID  uuid.UUID                  `db:"customer_id" json:"customerID"`
	Status      order_behavior.OrderStatus `db:"status" json:"status"`
	TotalAmount float64                    `db:"total_amount" json:"totalAmount"`
	Currency    string                     `db:"currency" json:"currency"`
	CreatedAt   sql.NullTime               `db:"created_at" json:"createdAt"`
	UpdatedAt   sql.NullTime               `db:"updated_at" json:"UpdatedAt"`
}

func (o *order) InsertOrder(ctx context.Context, order Order) error {
	_, err := o.database.ExecContext(ctx, "INSERT INTO orders (id, customer_id, status, total_amount, currency, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?, ?)",
		order.OrderID, order.CustomerID, order.Status, order.TotalAmount, order.Currency, order.CreatedAt, order.UpdatedAt)
	if err != nil {
		fmt.Printf("erro ao inserir os dados no banco, erro: %v", err)
		return err
	}

	err = o.redis.Set(ctx, "order_"+order.OrderID.String(), nil, 10*time.Minute)
	if err != nil {
		fmt.Printf("erro ao inserir no cache, erro: %v", err)
		return err
	}

	return nil
}
