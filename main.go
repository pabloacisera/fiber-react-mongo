package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	app := fiber.New()

	//middleware

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Use(logger.New())

	app.Static("", "./client/dist") //cada vez que realicemos un cambio en el frontend debe buildear de nuevo!

	app.Get("/works", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{ //es map es un struct
			"data": "lista de tareas",
		})
	})

	app.Listen(":3030")
	fmt.Println(">>> >>> Servidor corriendo en puerto 3030 >>> >>>")
}
