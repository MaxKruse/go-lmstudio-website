package requestdtos

import "github.com/openai/openai-go"

// @Description The Payload to use when asking the AI for help. You can provide the last params used to create a message chain
type CompletionRequest struct {
	Prompt     string                          `json:"prompt"`
	ParamsUsed *openai.ChatCompletionNewParams `json:"params_used,omitempty" swaggerignore:"true"`

	ParamsUsedJSON interface{} `json:"params_used,omitempty"`
}
