package config

import (
	"latihan2/middleware"
	"latihan2/route"

	"github.com/gofiber/fiber/v2"
)

func NewApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		},
	})

	app.Use(middleware.LoggerMiddleware)

	route.SetupRoutes(app)

	return app
}