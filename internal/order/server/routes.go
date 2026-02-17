package server

import (
	"log"
	"time"

	order_handler "github.com/andreparelho/order-api/internal/order/handler"
	order_service "github.com/andreparelho/order-api/internal/order/service"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/google/uuid"
)

func CreateRoute(orderService order_service.OrderService) *fiber.App {
	app := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	app.Use(recover.New())
	app.Use(requestid.New(requestid.Config{
		Header: "X-Request-ID",
		Generator: func() string {
			reqId, err := uuid.NewRandom()
			if err != nil {
				log.Fatal(err)
			}
			return reqId.String()
		},
	}))

	router := app.Group("/orders-api")
	router.Post("/orders", order_handler.OrderHandler(orderService))

	return app
}
