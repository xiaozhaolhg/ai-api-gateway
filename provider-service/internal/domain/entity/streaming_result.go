package entity

// StreamingResult encapsulates the result of a streaming transformation operation.
// It contains the transformed data, accumulated token counts, and a flag indicating
// whether this is the final chunk in the stream.
type StreamingResult struct {
	// TransformedData contains the transformed response data in OpenAI-compatible format.
	// For SSE streams, this is a single SSE chunk ready to be sent to the client.
	TransformedData []byte

	// TokenCounts contains the accumulated token counts for the stream.
	// During intermediate chunks, only AccumulatedTokens may be populated.
	// For the final chunk, PromptTokens and CompletionTokens will be set.
	TokenCounts TokenCounts

	// IsFinal indicates whether this is the final chunk in the stream.
	// When true, the stream has completed and all token counts are final.
	IsFinal bool
}
