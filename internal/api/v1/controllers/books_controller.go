package controllers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/maxkruse/go-lmstudio-website/internal/models/dtos"
	requestdtos "github.com/maxkruse/go-lmstudio-website/internal/models/dtos/request_dtos"
	"github.com/maxkruse/go-lmstudio-website/internal/service/book_service"
	"github.com/maxkruse/go-lmstudio-website/internal/utils/converters"
)

// @Summary Get all books
// @Description Gets all non-deleted books
// @Tags books
// @Produce json
// @Success 200 {array} dtos.Book
// @Failure	500	{object} dtos.ErrorResponse
// @Router /books [get]
func GetBooks(c echo.Context) error {
	books, err := book_service.Get()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, dtos.ErrorResponse{Error: err.Error()})
	}

	bookDtos := converters.BookEntityToDtoSlice(books)

	return c.JSON(http.StatusOK, bookDtos)
}

// @Summary Get a book by id
// @Description Gets a single book my id if possible
// @Tags books
// @Produce json
// @Param	id	path	integer	true "Book id"
// @Success 200 {object} dtos.Book
// @Failure	404	{object}	dtos.ErrorResponse
// @Router /books/{id} [get]
func GetBookById(c echo.Context) error {
	id := c.Param("id")
	// parse the id as an int32
	idInt, err := convertStringToInt32(id)
	if err != nil {
		return err
	}

	book, err := book_service.GetById(int32(idInt))

	if err != nil {
		return c.JSON(http.StatusNotFound, dtos.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, converters.BookEntityToDto(book))
}

// @Summary Create a book
// @Description Creates a new book
// @Tags books
// @Accept json
// @Produce json
// @Param	book	body	requestdtos.CreateBookRequest	true "Book to create"
// @Success 201 {object} dtos.Book
// @Failure	500	{object} dtos.ErrorResponse
// @Router /books [post]
func CreateBook(c echo.Context) error {
	var book requestdtos.CreateBookRequest
	if err := c.Bind(&book); err != nil {
		return err
	}

	newBook, err := book_service.Create(book)
	if err != nil {
		return c.JSON(http.StatusNotModified, dtos.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusCreated, converters.BookEntityToDto(newBook))
}

func UpdateBook(c echo.Context) error {
	var book dtos.Book
	if err := c.Bind(&book); err != nil {
		return err
	}
	id := c.Param("id")

	if err := book_service.Update(converters.BookDtoToEntity(book)); err != nil {
		return err
	}

	idInt, err := convertStringToInt32(id)
	if err != nil {
		return err
	}

	bookEntity, err := book_service.GetById(idInt)

	if err != nil {
		return c.JSON(http.StatusNotFound, dtos.ErrorResponse{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, converters.BookEntityToDto(bookEntity))
}

func DeleteBook(c echo.Context) error {
	id := c.Param("id")
	idInt, err := convertStringToInt32(id)
	if err != nil {
		return err
	}

	if err := book_service.Delete(idInt); err != nil {
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
