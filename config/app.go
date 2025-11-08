package config

import (
	"latihan2/middleware"
	"latihan2/route"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func NewApp() *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		},
	})

	app.Use(cors.New())

	app.Use(middleware.LoggerMiddleware)

	route.SetupRoutesPostgres(app)
	route.SetupRoutesMongo(app)

	return app
}
