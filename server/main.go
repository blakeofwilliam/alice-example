package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/blakeofwilliam/alice-example/server/handlers"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Save original os.Stdout
	originalStdout := os.Stdout

	// Create a buffer to hold the output
	var buf bytes.Buffer

	// Redirect os.Stdout to the buffer
	os.Stdout = os.NewFile(0, "/dev/null")
	defer func() {
		// Restore original os.Stdout
		os.Stdout = originalStdout
	}()

	// Create a multiwriter to write to both buffer and os.Stdout
	mw := io.MultiWriter(&buf, os.Stdout)
	fmt.Fprintf(mw, "")

	// Use the buffer's content as needed
	fmt.Println(buf.String())

	err := godotenv.Load()
	if err != nil {
		log.Println("<<< Unable to load .env file. If you're running this in production, this is not a problem.")
	}

	app := fiber.New()

	v1 := app.Group("/v1")

	v1.Get("/:operation", handlers.Operation)

	port := os.Getenv("SERVER_PORT")

	app.Listen(":" + port)
}