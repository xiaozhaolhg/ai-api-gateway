package port

// ProviderAdapter defines the interface for LLM provider adapters
type ProviderAdapter interface {
	// TransformRequest transforms the request to provider-specific format
	TransformRequest(request []byte, headers map[string]string) ([]byte, map[string]string, error)

	// TransformResponse transforms the response back to OpenAI format
	TransformResponse(response []byte) ([]byte, error)

	// CountTokens counts tokens in the request/response
	CountTokens(request []byte, response []byte) (int64, int64, error)
}
