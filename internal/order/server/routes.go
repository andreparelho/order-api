package server

import (
	"time"

	order_service "github.com/andreparelho/order-api/internal/order/service"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
)

func CreateRoute(orderService order_service.OrderService) *fiber.App {
	app := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})

	app.Use(recover.New())
	app.Use(requestid.New())

	router := app.Group("/orders-api")
	router.Post("/orders", orderService.CreateOrder())

	return app
}
