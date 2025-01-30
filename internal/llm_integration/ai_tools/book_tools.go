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
