package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/maxkruse/go-lmstudio-website/docs"
)

// @title LM Studio Website API
// @version 1.0
// @description This is a sample server for a bookstore server with included LLM functionality.

// @contact.name Maximilian Kruse
// @contact.url https://github.com/maxkruse

// @license.name No License

// @host bookstore.mkruse.xyz
// @BasePath /api/v1
func main() {
	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":3000"))
}
