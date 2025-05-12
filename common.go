package openai

// common.go defines common types used throughout the OpenAI API.

// Usage Represents the total token usage per request to OpenAI.
type Usage struct {
	PromptTokens            int                      `json:"prompt_tokens"`
	CompletionTokens        int                      `json:"completion_tokens"`
	TotalTokens             int                      `json:"total_tokens"`
	PromptTokensDetails     *PromptTokensDetails     `json:"prompt_tokens_details"`
	CompletionTokensDetails *CompletionTokensDetails `json:"completion_tokens_details"`
	InputTokens             int                      `json:"input_tokens"`
	OutputTokens            int                      `json:"output_tokens"`
	InputTokensDetails      InputTokensDetails       `json:"input_tokens_details"`
}

// CompletionTokensDetails Breakdown of tokens used in a completion.
type CompletionTokensDetails struct {
	AudioTokens          int `json:"audio_tokens"`
	ReasoningTokens      int `json:"reasoning_tokens"`
	CachedTokens         int `json:"cached_tokens"`
	CachedTokensInternal int `json:"cached_tokens_internal"`
	TextTokens           int `json:"text_tokens"`
	ImageTokens          int `json:"image_tokens"`
}

// PromptTokensDetails Breakdown of tokens used in the prompt.
type PromptTokensDetails struct {
	AudioTokens     int `json:"audio_tokens"`
	CachedTokens    int `json:"cached_tokens"`
	ReasoningTokens int `json:"reasoning_tokens"`
	TextTokens      int `json:"text_tokens"`
}

type InputTokensDetails struct {
	TextTokens  int `json:"text_tokens"`
	ImageTokens int `json:"image_tokens"`
}
