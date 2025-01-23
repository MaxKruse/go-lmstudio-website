package v1

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/maxkruse/go-lmstudio-website/internal/models/dtos"
	"github.com/maxkruse/go-lmstudio-website/internal/service"
	"github.com/maxkruse/go-lmstudio-website/internal/utils/converters"
)

// @Summary Get all books
// @Description Gets all non-deleted books
// @Tags books
// @Produce json
// @Success 200 {dtos.Book} Response
// @Router /path [get put post delete patch]
func GetBooks(c echo.Context) error {
	books := service.GetBooks()

	bookDtos := converters.BookEntityToDtoSlice(books)

	return c.JSON(http.StatusOK, bookDtos)
}

func GetBookById(c echo.Context) error {
	id := c.Param("id")
	// parse the id as an int32
	idInt, err := convertStringToInt32(id)
	if err != nil {
		return err
	}

	book := service.GetBookById(int32(idInt))
	return c.JSON(http.StatusOK, book)
}

func CreateBook(c echo.Context) error {
	var book dtos.Book
	if err := c.Bind(&book); err != nil {
		return err
	}

	if err := service.CreateBook(converters.BookDtoToEntity(book)); err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, book)
}

func UpdateBook(c echo.Context) error {
	var book dtos.Book
	if err := c.Bind(&book); err != nil {
		return err
	}
	id := c.Param("id")

	if err := service.UpdateBook(converters.BookDtoToEntity(book)); err != nil {
		return err
	}

	idInt, err := convertStringToInt32(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, converters.BookEntityToDto(service.GetBookById(idInt)))
}

func DeleteBook(c echo.Context) error {
	id := c.Param("id")
	idInt, err := convertStringToInt32(id)
	if err != nil {
		return err
	}

	if err := service.DeleteBook(idInt); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, nil)
}

func convertStringToInt32(s string) (int32, error) {
	i, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(i), nil
}
