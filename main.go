package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()

	// Variables de entorno
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))
	app.Use(logger.New())

	// Ruta para archivos estáticos
	//app.Static("/", "./client/dist")

	// Ruta API
	app.Get("/api/works", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"data": "lista de tareas",
		})
	})

	// Iniciar el servidor
	app.Listen(":" + port)
}
