package dto

type ChatCompletionRequestDTO struct {
	Model       string                   `json:"model" binding:"required"`
	Messages    []MessageDTO             `json:"messages" binding:"required,min=1"`
	Stream      bool                     `json:"stream"`
	Temperature float64                  `json:"temperature"`
	MaxTokens   int                      `json:"max_tokens"`
	TopP        float64                  `json:"top_p"`
}

type MessageDTO struct {
	Role    string `json:"role" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type ChatCompletionResponseDTO struct {
	ID      string        `json:"id"`
	Created int64         `json:"created"`
	Model   string        `json:"model"`
	Choices []ChoiceDTO   `json:"choices"`
	Usage   UsageDTO      `json:"usage"`
}

type ChoiceDTO struct {
	Index        int         `json:"index"`
	Message      MessageDTO  `json:"message"`
	FinishReason string      `json:"finish_reason"`
	Delta        *MessageDTO `json:"delta,omitempty"`
}

type UsageDTO struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

type ModelListDTO struct {
	Object string     `json:"object"`
	Data   []ModelDTO `json:"data"`
}

type ModelDTO struct {
	ID       string `json:"id"`
	Object   string `json:"object"`
	Created  int64  `json:"created"`
	OwnedBy  string `json:"owned_by"`
	Provider string `json:"provider"`
}