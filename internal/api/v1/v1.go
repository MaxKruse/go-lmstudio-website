package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/maxkruse/go-lmstudio-website/internal/api/v1/controllers"
)

func RegisterRoutes(e *echo.Echo) {
	books_group := e.Group("/api/v1/books")

	books_group.GET("", controllers.GetBooks)
	books_group.GET("/:id", controllers.GetBookById)
	books_group.POST("", controllers.CreateBook)
	books_group.PUT("/:id", controllers.UpdateBook)
	books_group.DELETE("/:id", controllers.DeleteBook)

	ai_group := e.Group("/api/v1/ai")

	ai_group.POST("/completion", controllers.AiChatCompletion)
}
