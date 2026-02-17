package order_repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	order_behavior "github.com/andreparelho/order-api/internal/order/behavior"
	"github.com/andreparelho/order-api/pkg/redis"
	"github.com/google/uuid"
)

type OrderRepository interface {
	InsertOrder(ctx context.Context, order Order, xRequestId string) (bool, error)
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

type OrderRedisInfo struct {
	Status order_behavior.OrderStatus `json:"status"`
}

func (o *order) InsertOrder(ctx context.Context, order Order, xRequestId string) (bool, error) {
	redisKey := fmt.Sprintf("order:req:id_%s", xRequestId)

	err := o.redis.Get(ctx, redisKey)
	if err != nil {
		fmt.Printf("WARN: nao foi possivel buscar o dado %s do cache. Criando dado na base de dados.", redisKey)
	} else {
		fmt.Printf("INFO: dado encotrado no cache.")
		return true, nil
	}

	_, err = o.database.ExecContext(ctx, "INSERT INTO orders (id, customer_id, status, total_amount, currency, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?, ?)",
		order.OrderID, order.CustomerID, order.Status, order.TotalAmount, order.Currency, order.CreatedAt, order.UpdatedAt)
	if err != nil {
		fmt.Printf("erro ao inserir os dados no banco, erro: %v", err)
		return false, err
	}

	orderInfo := OrderRedisInfo{
		Status: order_behavior.OrderStatusCreated,
	}

	orderInfoMarsh, err := json.Marshal(&orderInfo)
	if err != nil {
		fmt.Printf("erro ao realizar marshal, erro: %v", err)
		return false, err
	}

	err = o.redis.Set(ctx, redisKey, orderInfoMarsh, 10*time.Minute)
	if err != nil {
		fmt.Printf("erro ao inserir no cache, erro: %v", err)
		return false, err
	}

	return false, nil
}
