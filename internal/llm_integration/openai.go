package llm_integration

import (
	"context"
	"encoding/json"
	"log"
	"os"

	aitools "github.com/maxkruse/go-lmstudio-website/internal/llm_integration/ai_tools"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type AIClient struct {
	client         *openai.Client
	model          string
	availableTools []openai.ChatCompletionToolParam
}

func NewClient() AIClient {
	BASE_URL := os.Getenv("LM_STUDIO_HOST")
	baseUrlOption := option.WithBaseURL(BASE_URL)
	API_KEY := os.Getenv("LM_STUDIO_API_KEY")
	apiKeyOption := option.WithAPIKey(API_KEY)
	modelChoice := os.Getenv("LM_STUDIO_MODEL")

	client := openai.NewClient(baseUrlOption, apiKeyOption)

	var aiClient AIClient
	aiClient.client = client
	aiClient.model = modelChoice

	// get all available tools

	aiClient.addAllTools()

	return aiClient
}

func (ai *AIClient) addAllTools() {
	ai.availableTools = append(ai.availableTools, aitools.GetWeatherTool()...)
}

func (ai *AIClient) GetCompletion(ctx context.Context, prompt string) (string, error) {
	// make a completion with the entire prompt that we know

	params := openai.ChatCompletionNewParams{
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a highly capable and collaborative AI assistant with access to tools. Your purpose is to provide accurate, detailed, and context-aware responses to user queries. When you encounter a task that can benefit from tool usage (e.g., executing code, retrieving information, generating images, or processing data), you must use the appropriate tool. After using a tool, clearly explain the results or actions to the user. If the task requires multiple steps, plan your approach and communicate progress effectively. Always consider the context of the user's request and prioritize relevance, clarity, and precision in your responses. When unsure about a specific need, ask clarifying questions before proceeding. Avoid using tools unnecessarily and ensure any generated output is actionable and aligns with the user's goals."),
			openai.UserMessage(prompt),
		}),
		Tools:      openai.F(ai.availableTools),
		Seed:       openai.Int(0),
		Model:      openai.String(ai.model),
		ToolChoice: openai.F(openai.ChatCompletionToolChoiceOptionUnionParam(openai.ChatCompletionToolChoiceOptionAutoAuto)),
	}

	completion, err := ai.client.Chat.Completions.New(ctx, params)

	if err != nil {
		return "", err
	}

	toolCalls := completion.Choices[0].Message.ToolCalls
	if len(toolCalls) == 0 {
		log.Println("No tool calls found")
		return completion.Choices[0].Message.Content, nil
	}

	params.Messages.Value = append(params.Messages.Value, completion.Choices[0].Message)

	for _, toolCall := range toolCalls {
		if toolCall.Function.Name == "get_weather" {
			// get args if possible
			var args map[string]interface{}
			err := json.Unmarshal([]byte(toolCall.Function.Arguments), &args)

			if err != nil {
				log.Println("Error unmarshalling arguments:", err)
				continue
			}

			format, ok := args["format"].(string)
			if !ok {
				format = aitools.TEMP_FORMAT_CELSIUS
			}

			weatherResponse := aitools.GetWeatherFunc(format)

			params.Messages.Value = append(params.Messages.Value, openai.SystemMessage("You only respond to the question, you do not ask any followup questions or do smalltalk. If you cannot find the answer, respond in a way that asks the user to provide more information."))
			params.Messages.Value = append(params.Messages.Value, openai.ToolMessage(toolCall.ID, weatherResponse))
		}
	}

	completion, err = ai.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return "", err
	}

	return completion.Choices[0].Message.Content, nil
}
