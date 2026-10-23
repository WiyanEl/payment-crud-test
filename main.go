package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"

	"payment-crud-test/auth"
	"payment-crud-test/config"

	"payment-crud-test/middleware"

	"payment-crud-test/payments"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("File .env tidak ditemukan, menggunakan environment sistem")
	}

	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}
	defer db.Close()

	app := fiber.New()

	app.Get("/api", func(c *fiber.Ctx) error {
		if err := db.Ping(); err != nil {
			return c.Status(503).JSON(fiber.Map{
				"success": false,
				"message": "Database unavailable",
			})
		}

		return c.JSON(fiber.Map{
			"success":  true,
			"message":  "Payment API is running",
			"database": "connected",
		})
	})

	app.Post("/api/auth/register", auth.Register(db))
	app.Post("/api/auth/login", auth.Login(db))

	app.Get("/api/me", middleware.AuthRequired(), func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"success": true,
			"user_id": c.Locals("userID"),
		})
	})

	app.Get("/api/payments", middleware.AuthRequired(), payments.List(db))
	app.Get("/api/payments/:id", middleware.AuthRequired(), payments.GetByID(db))
	app.Post("/api/payments", middleware.AuthRequired(), payments.Create(db))
	app.Put("/api/payments/:id", middleware.AuthRequired(), payments.Update(db))
	app.Delete("/api/payments/:id", middleware.AuthRequired(), payments.Delete(db))

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(app.Listen(":" + port))
}