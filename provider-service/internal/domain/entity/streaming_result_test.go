package entity

import (
	"testing"
)

func TestStreamingResult_Initialization(t *testing.T) {
	// Test default initialization
	sr := StreamingResult{
		TransformedData: []byte("test data"),
		TokenCounts: TokenCounts{
			PromptTokens:      10,
			CompletionTokens:  20,
			AccumulatedTokens: 30,
		},
		IsFinal: false,
	}

	if string(sr.TransformedData) != "test data" {
		t.Errorf("Expected TransformedData to be 'test data', got %s", string(sr.TransformedData))
	}

	if sr.TokenCounts.PromptTokens != 10 {
		t.Errorf("Expected PromptTokens to be 10, got %d", sr.TokenCounts.PromptTokens)
	}

	if sr.TokenCounts.CompletionTokens != 20 {
		t.Errorf("Expected CompletionTokens to be 20, got %d", sr.TokenCounts.CompletionTokens)
	}

	if sr.TokenCounts.AccumulatedTokens != 30 {
		t.Errorf("Expected AccumulatedTokens to be 30, got %d", sr.TokenCounts.AccumulatedTokens)
	}

	if sr.IsFinal != false {
		t.Errorf("Expected IsFinal to be false, got %v", sr.IsFinal)
	}
}

func TestStreamingResult_FinalChunk(t *testing.T) {
	// Test final chunk configuration
	sr := StreamingResult{
		TransformedData: []byte("[DONE]"),
		TokenCounts: TokenCounts{
			PromptTokens:      15,
			CompletionTokens:  25,
			AccumulatedTokens: 40,
		},
		IsFinal: true,
	}

	if !sr.IsFinal {
		t.Error("Expected IsFinal to be true for final chunk")
	}

	total := sr.TokenCounts.Total()
	if total != 40 {
		t.Errorf("Expected Total() to be 40, got %d", total)
	}
}

func TestStreamingResult_EmptyData(t *testing.T) {
	// Test with empty transformed data
	sr := StreamingResult{
		TransformedData: []byte{},
		TokenCounts:     TokenCounts{},
		IsFinal:         true,
	}

	if len(sr.TransformedData) != 0 {
		t.Errorf("Expected empty TransformedData, got %d bytes", len(sr.TransformedData))
	}

	total := sr.TokenCounts.Total()
	if total != 0 {
		t.Errorf("Expected Total() to be 0, got %d", total)
	}
}
