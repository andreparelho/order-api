package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	order_service "github.com/andreparelho/order-api/internal/order/service"
	"github.com/andreparelho/order-api/pkg/config"
	"github.com/gofiber/fiber/v3"
)

type Server struct {
	App *fiber.App
}

func NewServer(cfg config.Configuration, orderService order_service.OrderService) (*Server, error) {
	app := CreateRoute(orderService)

	return &Server{
		App: app,
	}, nil
}

func (s *Server) Start(ctx context.Context, port string) error {
	go func() {
		<-ctx.Done()
		fmt.Print("shutting down http server")
		_ = s.App.Shutdown()
	}()

	fmt.Printf("starting server with port: %v", port)
	return s.App.Listen(":" + port)
}

func (s *Server) Shutdown() context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-ch
		cancel()
	}()

	return ctx
}
