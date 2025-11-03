package routes

import (

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewApp(client *mongo.Client) *fiber.App {
    app := fiber.New()
    
    app.Use(logger.New())

    app.Get("/", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{
            "message": "Welcome to the Alumni API",
            "success": true,
        })
    })

    return app
}