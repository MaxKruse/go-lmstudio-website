package aitools

import (
	"context"

	"github.com/maxkruse/go-lmstudio-website/internal/models/entities"
	"github.com/maxkruse/go-lmstudio-website/internal/service/book_service"
	"github.com/openai/openai-go"
)

func GetBookTools() []openai.ChatCompletionToolParam {

	toolData := []openai.ChatCompletionToolParam{
		{
			Type: openai.F(openai.ChatCompletionToolTypeFunction),
			Function: openai.F(openai.FunctionDefinitionParam{
				Name:        openai.String("get_books_by_price"),
				Description: openai.String("Gets books by price. Set the price to an insanely high number in case you want to see all the books."),
				Parameters: openai.F(openai.FunctionParameters{
					"type": "object",
					"properties": map[string]interface{}{
						"price": map[string]interface{}{
							"type": "number",
						},
					},
				}),
			}),
		},
	}

	return toolData
}

func GetBooksByPriceFunc(price float64) ([]entities.Book, error) {
	// make a database query to get all books below that price
	entities, err := book_service.GetBooksBelowPrice(context.Background(), price)

	return entities, err
}
