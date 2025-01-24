package dtos

import "github.com/openai/openai-go"

// @Description	Provides the entire context in case of back and forth communication, and also the last message as a handy accessor
type CompletionResult struct {
	ParamsUsed     openai.ChatCompletionNewParams `json:"params_used" swaggerignore:"true"`
	LastCompletion openai.ChatCompletionMessage   `json:"last_completion" swaggerignore:"true"`
	RawData        interface{}                    `json:"raw_data"`

	ParamsUsedJSON     interface{} `json:"-" example=Please visit the openai docs for 'ChatCompletionNewParams'"`
	LastCompletionJSON interface{} `json:"-" example=Please visit the openai docs for 'ChatCompletionMessage'"`
}
