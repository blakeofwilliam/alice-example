package main

import (
	"log"
	"os"

	"github.com/blakeofwilliam/alice-example/client/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("<<< Unable to load .env file. If you're running this in production, this is not a problem.")
	}

	app := fiber.New()

	v1 := app.Group("/v1")

	v1.Get("/generate", handlers.Generate)

	port := os.Getenv("CLIENT_PORT")

	app.Listen(":" + port)
}