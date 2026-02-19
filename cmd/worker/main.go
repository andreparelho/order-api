package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	order_consumer "github.com/andreparelho/order-api/internal/order/event"
	order_repository "github.com/andreparelho/order-api/internal/order/repository"
	payment_consumer "github.com/andreparelho/order-api/internal/payment/event"
	payment_event_repository "github.com/andreparelho/order-api/internal/payment/repository"
	"github.com/andreparelho/order-api/pkg/config"
	"github.com/andreparelho/order-api/pkg/rds"
	"github.com/andreparelho/order-api/pkg/redis"
	"github.com/andreparelho/order-api/pkg/sqs"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env não encontrado")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("erro ao carregar config: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	redisClient, err := redis.NewRedisClient(*cfg, ctx)
	if err != nil {
		log.Fatalf("erro redis: %v", err)
	}
	defer redisClient.Close()

	dbConn, err := rds.GetConnection(*cfg)
	if err != nil {
		log.Fatalf("erro rds: %v", err)
	}
	defer dbConn.Close()

	sqsClient := sqs.NewSQSClient(ctx, *cfg)

	orderEventRepository := order_repository.NewOrderEventRepository(sqsClient)
	orderRepository := order_repository.NewOrderRepository(dbConn, redisClient)

	paymentEventRepository := payment_event_repository.NewPaymentEventRepository(sqsClient)

	orderConsumer := order_consumer.NewOrderConsumer(*cfg, orderEventRepository, orderRepository)
	paymentConsumer := payment_consumer.NewPaymentConsumer(*cfg, paymentEventRepository)

	go orderConsumer.StartConsumer(ctx)
	go paymentConsumer.StartConsumer(ctx)

	fmt.Println("[INFO]: workers iniciados com sucesso")

	<-sig
	fmt.Println("[INFO]: shutdown signal recebido")

	cancel()
	fmt.Println("[INFO]: worker encerrado com sucesso")
}
