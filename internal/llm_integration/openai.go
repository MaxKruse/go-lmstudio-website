package llm_integration

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	aitools "github.com/maxkruse/go-lmstudio-website/internal/llm_integration/ai_tools"
	"github.com/maxkruse/go-lmstudio-website/internal/models/dtos"
	"github.com/maxkruse/go-lmstudio-website/internal/utils"
	"github.com/maxkruse/go-lmstudio-website/internal/utils/converters"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/valkey-io/valkey-go"
)

type AIClient struct {
	client         *openai.Client
	model          string
	availableTools []openai.ChatCompletionToolParam
	valkeyClient   valkey.Client
}

var aiClient AIClient
var once sync.Once

const SYSTEM_PROMPT = `You are a Book Price Assistant. Do not engage in small talk. Follow these rules:

1. TOOL USAGE:
- Required parameters: make sure required parameters are always included.
- Using Multiple Tools: You are allowed to use more than 1 tool at a time and use their combined output to make statements.


2. USER INPUT:
- Ask for price range: Always ask the user to provide a price range (price_min and price_max)
- Ask for author names.


NEVER mention technical details about the tool's implementation. Always keep prices formatted as $XX.XX.`

func NewClient() AIClient {

	once.Do(func() {
		BASE_URL := os.Getenv("LM_STUDIO_HOST")
		baseUrlOption := option.WithBaseURL(BASE_URL)
		API_KEY := os.Getenv("LM_STUDIO_API_KEY")
		apiKeyOption := option.WithAPIKey(API_KEY)
		modelChoice := os.Getenv("LM_STUDIO_MODEL")

		if BASE_URL == "" {
			log.Fatal("LM_STUDIO_HOST environment variable not set")
		}
		if API_KEY == "" {
			log.Fatal("LM_STUDIO_API_KEY environment variable not set")
		}
		if API_KEY == "LM_STUDIO_MODEL" {
			log.Fatal("LM_STUDIO_MODEL environment variable not set")
		}

		client := openai.NewClient(baseUrlOption, apiKeyOption)

		aiClient.client = client
		aiClient.model = modelChoice

		// get all available tools

		aiClient.addAllTools()

		// valkey setup
		valkeyClient, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{"localhost:6379"}})
		if err != nil {
			log.Fatal(err)
		}

		aiClient.valkeyClient = valkeyClient
	})

	return aiClient
}

func (ai *AIClient) addAllTools() {
	ai.availableTools = append(ai.availableTools, aitools.GetBookTools()...)
}

func (ai *AIClient) GetCompletion(ctx context.Context, prompt string, valkey_Key string) (*dtos.CompletionResult, error) {
	// make a completion with the entire prompt that we know

	var completionResult dtos.CompletionResult
	var params openai.ChatCompletionNewParams

	// check in valkey if we have data for this key
	jsonByteData, err := ai.valkeyClient.Do(context.Background(), ai.valkeyClient.B().Get().Key("ai-completions:"+valkey_Key).Build()).AsBytes()
	if err != nil {
		// no data reused
		log.Println("No previous ai completion data found for key: " + valkey_Key)
		params = openai.ChatCompletionNewParams{
			Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(SYSTEM_PROMPT),
			}),
			Tools:       openai.F(ai.availableTools),
			Seed:        openai.Int(0),
			Model:       openai.String(ai.model),
			Temperature: openai.Float(0.6),
		}
	} else {
		// we have the data, so we json unmarshal it
		var raw converters.RawChatData

		if err := json.Unmarshal(jsonByteData, &raw); err != nil {
			log.Printf("Couldnt unmarshal json byte data '%s' from valkey: %v", string(jsonByteData), err)
		}
		log.Println("Reusing previous ai completion data for key: " + valkey_Key)

		// use the parsed data as our params
		params = raw.ConvertToChatCompletionNewParams(ai.availableTools)
	}

	// attach the prompt after loading
	params.Messages.Value = append(params.Messages.Value, openai.UserMessage(prompt))

	completion, err := ai.client.Chat.Completions.New(ctx, params)

	if err != nil {
		log.Println("Error in completion:", err)
		return nil, err
	}

	// defer that we save the data in valkey and add a new key to the completionresult
	defer func() {
		prefix := "ai-completions:"
		newKey := utils.RandomString("", 10)
		paramsBytes, err := json.MarshalIndent(params, "", "  ")
		if err != nil {
			log.Printf("Failed to marshal params: %v", err)
			return // Or handle error differently
		}
		// we set the key to expire after 24 hours
		ai.valkeyClient.Do(context.Background(), ai.valkeyClient.B().Set().Key(prefix+newKey).Value(string(paramsBytes)).ExSeconds(86400).Build())
		completionResult.Key = newKey
	}()

	params.Messages.Value = append(params.Messages.Value, completion.Choices[0].Message)

	toolCalls := completion.Choices[0].Message.ToolCalls
	completionResult.Response = completion.Choices[0].Message.Content

	if len(toolCalls) == 0 {
		log.Println("No tool calls found")
		return &completionResult, nil
	}

	for _, toolCall := range toolCalls {
		var data interface{}
		var err error
		switch toolCall.Function.Name {
		case "get_books_by_price":
			data, err = handleGetBooksByPrice(toolCall)
		case "get_books_by_author":
			data, err = handleGetBooksByAuthor(toolCall)
		case "create_book":
			data, err = handleCreateBook(toolCall)
		default:
			err = fmt.Errorf("unknown tool: %s", toolCall.Function.Name)
		}

		// if the error is not nil, we tell the LLM that it failed
		if err != nil {
			params.Messages.Value = append(params.Messages.Value, openai.ToolMessage(
				toolCall.ID,
				fmt.Sprintf(`{"error": "%s"}`, err.Error()),
			))
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
	completionResult.Response = completion.Choices[0].Message.Content

	return &completionResult, nil
}
