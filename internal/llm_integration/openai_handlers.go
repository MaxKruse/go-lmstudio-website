package llm_integration

import (
	"encoding/json"
	"errors"
	"log"

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
