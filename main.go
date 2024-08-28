package main

import (
	"context"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/randomdev/go-fiber-react-vite-mongodb/models"
)

func main() {
	app := fiber.New()

	// Connect to Docker MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer client.Disconnect(context.TODO()) // Ensure disconnection

	coll := client.Database("gomongodb").Collection("works")

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Permite todos los orígenes
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	app.Use(logger.New())

	// Ruta para archivos estáticos
	app.Static("/", "../cliente/dist")

	// Ruta API
	app.Get("/api/works", func(c *fiber.Ctx) error {
		var tasks []models.Task

		// Realiza la consulta para obtener todos los documentos
		cur, err := coll.Find(context.TODO(), bson.D{})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error al obtener tareas de la BD",
			})
		}
		defer cur.Close(context.TODO())

		// Itera sobre el cursor y decodifica cada documento en el slice de tareas
		if err := cur.All(context.TODO(), &tasks); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error al procesar los datos",
			})
		}

		// Devuelve la lista de tareas en formato JSON
		return c.JSON(tasks)
	})

	app.Post("/api/create", func(c *fiber.Ctx) error {
		var task models.Task

		// Imprime el cuerpo de la solicitud
		bodyBytes := c.Body()
		log.Printf("Received body: %s", bodyBytes)

		if err := c.BodyParser(&task); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "No se encontró el body",
			})
		}

		task.Confirmed = false

		_, err := coll.InsertOne(context.TODO(), bson.D{
			{"title", task.Title},
			{"description", task.Description},
			{"deadline", task.Deadline},
			{"confirmed", task.Confirmed},
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error al insertar documento en BD",
			})
		}

		return c.JSON(fiber.Map{
			"message": "Tarea creada correctamente",
		})
	})

	app.Delete("/api/delete/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "ID no válido",
			})
		}
		_, err = coll.DeleteOne(context.TODO(), bson.D{{"_id", objId}})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error al eliminar documento de BD",
			})
		}
		return c.JSON(fiber.Map{
			"message": "Tarea eliminada correctamente",
		})
	})

	app.Put("/api/confirm/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		objId, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "ID no válido",
			})
		}
		update := bson.D{{
			"$set", bson.D{
				{"confirmed", true}},
		}}
		_, err = coll.UpdateOne(context.TODO(), bson.D{{"_id", objId}}, update)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error al confirmar la tarea de BD",
			})
		}
		return c.JSON(fiber.Map{
			"message": "Tarea confirmada correctamente",
		})
	})

	// Variables de entorno
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Iniciar el servidor
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
