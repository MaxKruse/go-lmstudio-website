package aitools

import (
	"github.com/openai/openai-go"
)

func GetWeatherTool() []openai.ChatCompletionToolParam {

	toolData := []openai.ChatCompletionToolParam{
		{
			Type: openai.F(openai.ChatCompletionToolTypeFunction),
			Function: openai.F(openai.FunctionDefinitionParam{
				Name:        openai.String("get_weather"),
				Description: openai.String("Gets the weather for a given location"),
				Parameters: openai.F(openai.FunctionParameters{
					"type": "object",
					"properties": map[string]interface{}{
						"format": map[string]interface{}{
							"type": "string",
							"enum": []string{TEMP_FORMAT_CELSIUS, TEMP_FORMAT_FAHRENHEIT},
						},
					},
				}),
			}),
		},
	}

	return toolData
}

const TEMP_FORMAT_CELSIUS = "celsius"
const TEMP_FORMAT_FAHRENHEIT = "fahrenheit"

func GetWeatherFunc(format string) string {
	if format == TEMP_FORMAT_CELSIUS {
		return "26.5°C"
	}

	return "74.4°F"
}
