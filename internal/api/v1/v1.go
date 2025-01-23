package v1

import "github.com/labstack/echo/v4"

func RegisterRoutes(e *echo.Echo) {
	booksgroup := e.Group("/api/v1/books")

	booksgroup.GET("", GetBooks)
	booksgroup.GET("/:id", GetBookById)
	booksgroup.POST("", CreateBook)
	booksgroup.PUT("/:id", UpdateBook)
	booksgroup.DELETE("/:id", DeleteBook)
}
