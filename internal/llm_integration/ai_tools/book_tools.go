package aitools

import (
	"context"
	"time"

	"github.com/maxkruse/go-lmstudio-website/internal/models/dtos"
	requestdtos "github.com/maxkruse/go-lmstudio-website/internal/models/dtos/request_dtos"
	"github.com/maxkruse/go-lmstudio-website/internal/models/entities"
	"github.com/maxkruse/go-lmstudio-website/internal/service/book_service"
	"github.com/maxkruse/go-lmstudio-website/internal/utils/converters"
	"github.com/openai/openai-go"
)

func GetBookTools() []openai.ChatCompletionToolParam {

	toolData := []openai.ChatCompletionToolParam{
		{
			Type: openai.F(openai.ChatCompletionToolTypeFunction),
			Function: openai.F(openai.FunctionDefinitionParam{
				Name:        openai.String("get_books_by_price"),
				Description: openai.String("Gets books by price. Set the arguments to either 0 for min or 1000000 for max if one of them is not needed."),
				Parameters: openai.F(openai.FunctionParameters{
					"type": "object",
					"properties": map[string]interface{}{
						"price_min": map[string]interface{}{
							"type": "number",
						},
						"price_max": map[string]interface{}{
							"type": "number",
						},
					},
					"required": []string{"price_min", "price_max"},
				}),
			}),
		},
		{
			Type: openai.F(openai.ChatCompletionToolTypeFunction),
			Function: openai.F(openai.FunctionDefinitionParam{
				Name:        openai.String("get_books_by_author"),
				Description: openai.String("Gets books by author. The author is case sensitive and needs to be spelled correctly."),
				Parameters: openai.F(openai.FunctionParameters{
					"type": "object",
					"properties": map[string]interface{}{
						"author": map[string]interface{}{
							"type": "string",
						},
					},
					"required": []string{"author"},
				}),
			}),
		},
		{
			Type: openai.F(openai.ChatCompletionToolTypeFunction),
			Function: openai.F(openai.FunctionDefinitionParam{
				Name:        openai.String("create_book"),
				Description: openai.String("Creates a new book. All information needs to be present to do so. The user will present the price in $ or €, so we strip that and convert it to a float32."),
				Parameters: openai.F(openai.FunctionParameters{
					"type": "object",
					"properties": map[string]interface{}{
						"title": map[string]interface{}{
							"type": "string",
						},
						"author": map[string]interface{}{
							"type": "string",
						},
						"price": map[string]interface{}{
							"type": "number",
						},
						"img_url": map[string]interface{}{
							"type": "string",
						},
						"isbn": map[string]interface{}{
							"type": "string",
						},
						"published_date": map[string]interface{}{
							"type": "string",
						},
						"description": map[string]interface{}{
							"type": "string",
						},
					},
					"required": []string{"title", "author", "price", "img_url", "isbn", "published_date", "description"},
				}),
			}),
		},
	}

	return toolData
}

func GetBooksByPriceFunc(price_min float64, price_max float64) ([]entities.Book, error) {
	// make a database query to get all books below that price
	entities, err := book_service.GetBooksBetweenPrice(context.Background(), price_min, price_max)

	return entities, err
}

func GetBooksByAuthorFunc(author string) ([]entities.Book, error) {
	entities, err := book_service.GetBooksByAuthor(context.Background(), author)
	return entities, err
}

func CreateBookFunc(title string, author string, price float32, imgUrl string, isbn string, published time.Time, desc string) (dtos.Book, error) {
	var createBook requestdtos.CreateBookRequest

	createBook.Author = author
	createBook.Title = title
	createBook.Price = price
	createBook.ImageUrl = imgUrl // TODO: add image url functionality
	createBook.Isbn = isbn       // TODO: add isbn functionality
	createBook.PublishedDate = published.Format("2006-01-02")
	createBook.Description = desc

	newBook, err := book_service.Create(createBook)
	return converters.BookEntityToDto(newBook), err
}
