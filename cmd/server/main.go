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

// @contact.name None
// @contact.url None
// @contact.email None

// @license.name None
// @license.url None

// @host bookstore.mkruse.xyz
// @BasePath /v1
func main() {
	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":3000"))
}
