package llm_integration

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	aitools "github.com/maxkruse/go-lmstudio-website/internal/llm_integration/ai_tools"
	"github.com/openai/openai-go"
)

func handleGetBooksByPrice(toolCall openai.ChatCompletionMessageToolCall) (interface{}, error) {
	// Get args if possible
	var args map[string]interface{}
	err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
	if err != nil {
		log.Println("Error unmarshalling arguments:", err)
		return "", err
	}

	priceMin, ok := args["price_min"].(float64)
	if !ok {
		return "", errors.New("price_min is missing or not a float64")
	}

	priceMax, ok := args["price_max"].(float64)
	if !ok {
		return "", errors.New("price_max is missing or not a float64")
	}

	return aitools.GetBooksByPriceFunc(priceMin, priceMax)
}

func handleGetBooksByAuthor(toolCall openai.ChatCompletionMessageToolCall) (interface{}, error) {
	var args map[string]interface{}
	err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
	if err != nil {
		log.Println("Error unmarshalling arguments:", err)
		return "", err
	}

	author, ok := args["author"].(string)
	if !ok {
		return "", errors.New("author is missing or not a string")
	}

	return aitools.GetBooksByAuthorFunc(author)
}

func handleCreateBook(toolCall openai.ChatCompletionMessageToolCall) (interface{}, error) {
	var args map[string]interface{}
	err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)
	if err != nil {
		log.Println("Error unmarshalling arguments:", err)
		return "", err
	}

	title, ok := args["title"].(string)
	if !ok {
		return "", errors.New("title is missing or not a string")
	}
	author, ok := args["author"].(string)
	if !ok {
		return "", errors.New("author is missing or not a string")
	}
	price, ok := args["price"].(float64)
	if !ok {
		return "", errors.New("price is missing or not a float64")
	}
	imgUrl, ok := args["img_url"].(string)
	if !ok {
		return "", errors.New("img_url is missing or not a string")
	}
	isbn, ok := args["isbn"].(string)
	if !ok {
		return "", errors.New("isbn is missing or not a string")
	}
	publishedDateStr, ok := args["published_date"].(string)
	if !ok {
		return "", errors.New("published_date is missing or not a string")
	}
	publishedDate, err := time.Parse("2006-01-02", publishedDateStr)
	if err != nil {
		return "", errors.New("published_date is not valid")
	}
	desc, ok := args["description"].(string)
	if !ok {
		return "", errors.New("description is missing or not a string")
	}

	return aitools.CreateBookFunc(title, author, float32(price), imgUrl, isbn, publishedDate, desc)
}
