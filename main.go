package main

import (
	"context"
	"fmt"
	"log"

	order_repository "github.com/andreparelho/order-api/internal/order/repository"
	"github.com/andreparelho/order-api/internal/order/server"
	order_service "github.com/andreparelho/order-api/internal/order/service"
	"github.com/andreparelho/order-api/pkg/config"
	"github.com/andreparelho/order-api/pkg/rds"
	"github.com/andreparelho/order-api/pkg/redis"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Print(" .env não encontrado")
		log.Fatal(err)
	}

	config, err := config.Load()
	if err != nil {
		fmt.Printf("erro ao carregar as configuracoes do projeto. erro: %v", err)
		log.Fatal(err)
	}

	ctx := context.Background()

	redis, err := redis.NewRedisClient(*config, ctx)
	if err != nil {
		fmt.Printf("erro ao conectar com redis. erro: %v", err)
		log.Fatal(err)
	}
	defer redis.Close()

	dbConn, err := rds.GetConnection(*config)
	if err != nil {
		fmt.Printf("erro ao buscar conexao com rds. erro: %v", err)
		log.Fatal(err)
	}
	defer dbConn.Close()

	orderRepository := order_repository.NewOrderRepository(dbConn, redis)

	orderService := order_service.NewOrderService(orderRepository)
	server, err := server.NewServer(*config, orderService)

	if err := server.Start(ctx, config.Port); err != nil {
		fmt.Printf("erro ao inicializar a aplicacao. erro: %v", err)
		log.Fatal(err)
	}
}
