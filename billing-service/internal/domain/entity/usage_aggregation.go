package entity

// UsageAggregation represents aggregated usage statistics
type UsageAggregation struct {
	UserID           string  `json:"user_id"`
	GroupID          string  `json:"group_id"`
	ProviderID       string  `json:"provider_id"`
	Model            string  `json:"model"`
	TotalRequests    int64   `json:"total_requests"`
	TotalPromptTokens    int64   `json:"total_prompt_tokens"`
	TotalCompletionTokens int64 `json:"total_completion_tokens"`
	TotalCost        float64 `json:"total_cost"`
	StartDate        string  `json:"start_date"`
	EndDate          string  `json:"end_date"`
}
