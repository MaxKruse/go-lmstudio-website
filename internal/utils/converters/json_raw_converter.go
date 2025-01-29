package converters

import (
	"encoding/json"
	"fmt"

	"github.com/openai/openai-go"
)

type RawChatData struct {
	Messages    []Message `json:"messages"`
	Model       string    `json:"model"`
	Seed        int       `json:"seed"`
	Temperature float64   `json:"temperature"`
	Tools       []Tool    `json:"tools"`
}

type Message struct {
	Role       string      `json:"role"`
	Content    interface{} `json:"content,omitempty"` // Can be []ContentItem or string
	ToolCalls  []ToolCall  `json:"tool_calls,omitempty"`
	ToolCallID string      `json:"tool_call_id,omitempty"`
}

type ContentItem struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type ToolCall struct {
	ID       string           `json:"id"`
	Type     string           `json:"type"`
	Function ToolCallFunction `json:"function"`
}

type ToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type Tool struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

type ToolFunction struct {
	Description string     `json:"description"`
	Name        string     `json:"name"`
	Parameters  Parameters `json:"parameters"`
}

type Parameters struct {
	Type       string                       `json:"type"`
	Properties map[string]ParameterProperty `json:"properties"`
	Required   []string                     `json:"required"`
}

type ParameterProperty struct {
	Type string `json:"type"`
}

// Custom UnmarshalJSON for Message
func (m *Message) UnmarshalJSON(data []byte) error {
	type Alias Message
	aux := &struct {
		Content json.RawMessage `json:"content,omitempty"`
		*Alias
	}{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Content != nil {
		// Try to unmarshal as []ContentItem first
		var items []ContentItem
		if err := json.Unmarshal(aux.Content, &items); err == nil {
			m.Content = items
		} else {
			// Try to unmarshal as string
			var str string
			if err := json.Unmarshal(aux.Content, &str); err == nil {
				m.Content = str
			} else {
				return fmt.Errorf("failed to unmarshal content")
			}
		}
	}
	return nil
}

// Custom MarshalJSON for Message
func (m *Message) MarshalJSON() ([]byte, error) {
	type Alias Message
	var content interface{}

	switch v := m.Content.(type) {
	case []ContentItem:
		content = v
	case string:
		content = v
	default:
		content = nil
	}

	return json.Marshal(&struct {
		Content interface{} `json:"content,omitempty"`
		*Alias
	}{
		Content: content,
		Alias:   (*Alias)(m),
	})
}

func (r *RawChatData) ConvertToChatCompletionNewParams(availableTools []openai.ChatCompletionToolParam) openai.ChatCompletionNewParams {
	var res openai.ChatCompletionNewParams

	// we need to fill the following fields correctly:
	// - messages
	// - model
	// - temperature
	// - seed
	// - tools

	// messages
	var unions []openai.ChatCompletionMessageParamUnion

	// helper function that determines if the message content is a string, a list of ContentItems, or somethign else
	isString := func(content interface{}) bool {
		switch content.(type) {
		case string:
			return true
		case []ContentItem:
			return false
		default:
			panic("unknown content type")
		}
	}

	isContentItemArray := func(content interface{}) bool {
		switch content.(type) {
		case []ContentItem:
			return true
		case string:
			return false
		default:
			panic("unknown content type")
		}
	}

	for _, message := range r.Messages {

		// check if its a tool call, system, assistant or user message
		switch message.Role {
		case "system":
			if isString(message.Content) {
				unions = append(unions, openai.SystemMessage(message.Content.(string)))
			} else if isContentItemArray(message.Content) {
				for _, contentItem := range message.Content.([]ContentItem) {
					unions = append(unions, openai.SystemMessage(contentItem.Text))
				}
			}
		case "assistant":
			if isString(message.Content) {
				// its just text
				unions = append(unions, openai.AssistantMessage(message.Content.(string)))
			} else if isContentItemArray(message.Content) {
				// its a tool call
				// we ignore these, the data is already filled by the message itself
			}

		case "user":
			if isString(message.Content) {
				unions = append(unions, openai.UserMessage(message.Content.(string)))
			} else if isContentItemArray(message.Content) {
				for _, contentItem := range message.Content.([]ContentItem) {
					unions = append(unions, openai.UserMessage(contentItem.Text))
				}
			}
		case "tool":
			if isString(message.Content) {
				unions = append(unions, openai.ToolMessage(message.ToolCallID, message.Content.(string)))
			} else if isContentItemArray(message.Content) {
				for _, contentItem := range message.Content.([]ContentItem) {
					unions = append(unions, openai.ToolMessage(message.ToolCallID, contentItem.Text))
				}
			}
		}
	}
	res.Messages = openai.F(unions)

	res.Model = openai.String(r.Model)
	res.Temperature = openai.Float(r.Temperature)
	res.Seed = openai.Int(int64(r.Seed))
	res.Tools = openai.F(availableTools)

	return res
}
