package entity

// TokenCounts represents token usage information for a request/response.
// It is used for both streaming and non-streaming scenarios.
type TokenCounts struct {
	// PromptTokens is the number of tokens in the prompt/request.
	// For streaming, this is only populated in the final chunk.
	PromptTokens int64 `json:"prompt_tokens"`

	// CompletionTokens is the number of tokens in the completion/response.
	// For streaming, this is only populated in the final chunk.
	CompletionTokens int64 `json:"completion_tokens"`

	// AccumulatedTokens is the running total of tokens accumulated during streaming.
	// This is updated progressively during intermediate chunks and represents
	// the total tokens seen so far (prompt + completion).
	AccumulatedTokens int64 `json:"accumulated_tokens"`
}

// Total returns the total token count (prompt + completion).
// If AccumulatedTokens is set (streaming scenario), it returns that value.
// Otherwise returns PromptTokens + CompletionTokens.
func (t TokenCounts) Total() int64 {
	if t.AccumulatedTokens > 0 {
		return t.AccumulatedTokens
	}
	return t.PromptTokens + t.CompletionTokens
}
