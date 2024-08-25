package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	app := fiber.New()

	//connect to docker-mongo
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017/gomongodb"))

	if err != nil {
		panic(err)
	}

	coll := client.Database("gomongodb").Collection("works")

	_, err = coll.InsertOne(context.TODO(), bson.M{
		"task": "aprender a programar",
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Documento insertado exitosamente!")

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
	app.Static("/", "../cliente/dist")

	// Ruta API
	app.Get("/api/works", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"data": "lista de tareas",
		})
	})

	// Iniciar el servidor
	app.Listen(":" + port)
}
