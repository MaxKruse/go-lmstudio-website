package dtos

// @Description	Provides the entire context in case of back and forth communication, and also the last message as a handy accessor
type CompletionResult struct {
	Key      string      `json:"key"`
	Response string      `json:"Response"`
	RawData  interface{} `json:"raw_data"`
}
