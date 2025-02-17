package aitools

import (
	"time"

	"github.com/openai/openai-go"
)

// general tools for things like "whats the current date and time"

func GetUtilTools() []openai.ChatCompletionToolParam {
	toolData := []openai.ChatCompletionToolParam{
		{
			Type: openai.F(openai.ChatCompletionToolTypeFunction),
			Function: openai.F(openai.FunctionDefinitionParam{
				Name:        openai.String("get_current_date"),
				Description: openai.String("Gets the current date and time."),
				Parameters:  openai.F(openai.FunctionParameters{}),
			}),
		},
	}

	return toolData
}

func GetCurrentDateFunc() (string, error) {
	return time.Now().Format("2006-01-02 15:04:05"), nil
}
