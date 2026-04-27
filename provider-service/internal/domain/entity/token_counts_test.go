package entity

import (
	"testing"
)

func TestTokenCounts_Total_WithAccumulated(t *testing.T) {
	// When AccumulatedTokens is set, Total() should return that value
	tc := TokenCounts{
		PromptTokens:      10,
		CompletionTokens:  20,
		AccumulatedTokens: 35,
	}

	total := tc.Total()
	if total != 35 {
		t.Errorf("Expected Total() to return AccumulatedTokens (35), got %d", total)
	}
}

func TestTokenCounts_Total_WithoutAccumulated(t *testing.T) {
	// When AccumulatedTokens is 0, Total() should return PromptTokens + CompletionTokens
	tc := TokenCounts{
		PromptTokens:      10,
		CompletionTokens:  20,
		AccumulatedTokens: 0,
	}

	total := tc.Total()
	expected := int64(30)
	if total != expected {
		t.Errorf("Expected Total() to be %d, got %d", expected, total)
	}
}

func TestTokenCounts_FieldAccess(t *testing.T) {
	tc := TokenCounts{
		PromptTokens:      100,
		CompletionTokens:  200,
		AccumulatedTokens: 300,
	}

	if tc.PromptTokens != 100 {
		t.Errorf("Expected PromptTokens to be 100, got %d", tc.PromptTokens)
	}

	if tc.CompletionTokens != 200 {
		t.Errorf("Expected CompletionTokens to be 200, got %d", tc.CompletionTokens)
	}

	if tc.AccumulatedTokens != 300 {
		t.Errorf("Expected AccumulatedTokens to be 300, got %d", tc.AccumulatedTokens)
	}
}

func TestTokenCounts_ZeroValues(t *testing.T) {
	tc := TokenCounts{}

	if tc.PromptTokens != 0 {
		t.Errorf("Expected default PromptTokens to be 0, got %d", tc.PromptTokens)
	}

	if tc.CompletionTokens != 0 {
		t.Errorf("Expected default CompletionTokens to be 0, got %d", tc.CompletionTokens)
	}

	if tc.AccumulatedTokens != 0 {
		t.Errorf("Expected default AccumulatedTokens to be 0, got %d", tc.AccumulatedTokens)
	}

	total := tc.Total()
	if total != 0 {
		t.Errorf("Expected Total() to be 0 for zero values, got %d", total)
	}
}
