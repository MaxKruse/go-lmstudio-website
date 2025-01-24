package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/maxkruse/go-lmstudio-website/docs"
	v1 "github.com/maxkruse/go-lmstudio-website/internal/api/v1"
	"github.com/maxkruse/go-lmstudio-website/internal/db/migrations"
)

func init() {
	// load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

// @title LM Studio Website API
// @version 1.0
// @description This is a sample server for a bookstore server with included LLM functionality.

// @contact.name Maximilian Kruse
// @contact.url https://github.com/maxkruse

// @license.name No License

// @host localhost:3000
// @BasePath /api/v1
func main() {
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "3000"
	}

	env := os.Getenv("SERVER_ENV")
	if env == "" {
		env = "development"
	}

	fmt.Printf("Starting server on port %s in %s mode...\n", port, env)

	// in case anything goes boom, defer a recover
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Panic: %v", r)
		}
	}()

	// echo base router
	e := echo.New()

	// in development mode, use the swagger docs and allow CORS
	if env == "development" {
		e.Use(middleware.CORS())
		e.GET("/swagger/*", echoSwagger.WrapHandler)

		if err := migrations.MigrateDown(); err != nil {
			log.Fatal(err)
		}

		// also migrate up
		if err := migrations.MigrateUp(); err != nil {
			log.Fatal(err)
		}
	}

	v1.RegisterRoutes(e)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", port)))
}
