package main

import (
	"context"
	"log"

	"github.com/andreparelho/order-api/internal/order/server"
	order_service "github.com/andreparelho/order-api/internal/order/service"
	"github.com/andreparelho/order-api/pkg/config"
	"github.com/andreparelho/order-api/pkg/redis"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(" .env não encontrado")
		log.Fatal(err)
	}

	config, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	redis, err := redis.NewRedisClient(*config, ctx)

	orderService := order_service.NewOrderService(redis)
	server, err := server.NewServer(*config, orderService)

	if err := server.Start(ctx, config.Port); err != nil {
		log.Fatal(err)
	}
}
