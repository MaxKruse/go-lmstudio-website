package requestdtos

// @Description The Payload to use when asking the AI for help. You can provide the last params used to create a message chain
type CompletionRequest struct {
	Prompt string `json:"prompt"`
	Key    string `json:"key,omitempty"`
}
