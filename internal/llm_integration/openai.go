package llm_integration

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	aitools "github.com/maxkruse/go-lmstudio-website/internal/llm_integration/ai_tools"
	"github.com/maxkruse/go-lmstudio-website/internal/models/dtos"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type AIClient struct {
	client         *openai.Client
	model          string
	availableTools []openai.ChatCompletionToolParam
}

var aiClient AIClient

func NewClient() AIClient {

	if aiClient.client != nil {
		return aiClient
	}

	BASE_URL := os.Getenv("LM_STUDIO_HOST")
	baseUrlOption := option.WithBaseURL(BASE_URL)
	API_KEY := os.Getenv("LM_STUDIO_API_KEY")
	apiKeyOption := option.WithAPIKey(API_KEY)
	modelChoice := os.Getenv("LM_STUDIO_MODEL")

	client := openai.NewClient(baseUrlOption, apiKeyOption)

	aiClient.client = client
	aiClient.model = modelChoice

	// get all available tools

	aiClient.addAllTools()

	return aiClient
}

func (ai *AIClient) addAllTools() {
	ai.availableTools = append(ai.availableTools, aitools.GetBookTools()...)
}

func (ai *AIClient) GetCompletion(ctx context.Context, prompt string, paramsUsed *openai.ChatCompletionNewParams) (*dtos.CompletionResult, error) {
	// make a completion with the entire prompt that we know

	var completionResult dtos.CompletionResult
	var params openai.ChatCompletionNewParams

	if paramsUsed != nil {
		params = *paramsUsed
	} else {
		params = openai.ChatCompletionNewParams{
			Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage("You are a highly capable and collaborative AI assistant for a bookshop website with access to tools. Your purpose is to provide accurate, detailed, and context-aware responses to user queries. When you encounter a task that can benefit from tool usage (e.g., executing code, retrieving information, generating images, or processing data), you must use the appropriate tool. After using a tool, clearly explain the results or actions to the user. If the task requires multiple steps, plan your approach and communicate progress effectively. Always consider the context of the user's request and prioritize relevance, clarity, and precision in your responses. When unsure about a specific need, ask clarifying questions before proceeding. Ensure any output is actionable and aligns with the user's goals."),
				openai.UserMessage(prompt),
			}),
			Tools:       openai.F(ai.availableTools),
			Seed:        openai.Int(0),
			Model:       openai.String(ai.model),
			ToolChoice:  openai.F(openai.ChatCompletionToolChoiceOptionUnionParam(openai.ChatCompletionToolChoiceOptionAutoAuto)),
			Temperature: openai.Float(0.6),
		}
	}

	completion, err := ai.client.Chat.Completions.New(ctx, params)

	if err != nil {
		return nil, err
	}

	params.Messages.Value = append(params.Messages.Value, completion.Choices[0].Message)

	toolCalls := completion.Choices[0].Message.ToolCalls

	completionResult.ParamsUsed = params
	completionResult.LastCompletion = completion.Choices[0].Message

	if len(toolCalls) == 0 {
		log.Println("No tool calls found")
		return &completionResult, nil
	}

	params.Messages.Value = append(params.Messages.Value, completion.Choices[0].Message)
	// tell the model how to handle tools, since we need to use at least one now
	params.Messages.Value = append(params.Messages.Value, openai.SystemMessage(
		"You only respond to the question, you do not ask any followup questions or do smalltalk. If you cannot find the answer, respond in a way that asks the user to provide more information and tell them what error we encountered, in a way a non-technical person can understand. Behave like a salesperson, so if they ask for something that we dont have, give them a hint for other things we might have."))

	for _, toolCall := range toolCalls {
		var data interface{}
		var err error
		if toolCall.Function.Name == "get_books_by_price" {
			data, err = handleGetBooksByPrice(toolCall)
		}

		// if the error is not nil, we tell the LLM that it failed
		if err != nil {
			params.Messages.Value = append(params.Messages.Value, openai.SystemMessage(fmt.Sprintf("Error: %s", err.Error())))
			params.Messages.Value = append(params.Messages.Value, openai.ToolMessage(toolCall.ID, "Error: "+err.Error()))
			continue
		}

		byteResult, _ := json.Marshal(data)

		params.Messages.Value = append(params.Messages.Value, openai.ToolMessage(toolCall.ID, string(byteResult)))
		completionResult.RawData = data
	}

	completion, err = ai.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return nil, err
	}

	params.Messages.Value = append(params.Messages.Value, completion.Choices[0].Message)

	completionResult.ParamsUsed = params
	completionResult.LastCompletion = completion.Choices[0].Message

	return &completionResult, nil
}
