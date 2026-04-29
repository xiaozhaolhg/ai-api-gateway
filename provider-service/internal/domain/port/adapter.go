package port

import "github.com/ai-api-gateway/provider-service/internal/domain/entity"

// ProviderAdapter defines the interface for LLM provider adapters.
// It supports both streaming (SSE) and non-streaming response transformations.
type ProviderAdapter interface {
	// TransformRequest transforms the request to provider-specific format.
	// The request is in OpenAI-compatible format and is converted to the provider's native format.
	TransformRequest(request []byte, headers map[string]string) ([]byte, map[string]string, error)

	// TransformResponse transforms the response back to OpenAI format.
	//
	// For non-streaming (isStreaming=false):
	//   - response contains the complete response body
	//   - accumulatedTokens should be empty
	//   - returns (transformedData, tokenCounts, isFinal=true, error)
	//
	// For streaming (isStreaming=true):
	//   - response contains a single SSE chunk
	//   - accumulatedTokens contains the accumulated state from previous chunks
	//   - returns (transformedChunk, updatedTokenCounts, isFinal, error)
	//   - isFinal is true when the chunk contains the stream termination marker
	TransformResponse(response []byte, isStreaming bool, accumulatedTokens entity.TokenCounts) ([]byte, entity.TokenCounts, bool, error)

	// CountTokens counts tokens in the request/response.
	//
	// For non-streaming (isStreaming=false):
	//   - Returns the actual prompt and completion token counts
	//
	// For streaming (isStreaming=true):
	//   - Returns (0,0) for intermediate chunks
	//   - Returns actual counts for the final chunk
	//   - May estimate tokens if provider doesn't provide explicit counts
	CountTokens(request []byte, response []byte, isStreaming bool) (int64, int64, error)

	TestConnection(credentials string) error
}
