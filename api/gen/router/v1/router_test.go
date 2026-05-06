package routerv1

import (
	"testing"
	"google.golang.org/protobuf/proto"
)

func TestRoutingRule_NewFields(t *testing.T) {
	rule := &RoutingRule{
		Id:                 "rule-1",
		UserId:              "user-123",
		ModelPattern:        "ollama:*",
		ProviderId:          "ollama",
		Priority:            1,
		FallbackProviderIds: []string{"opencode_zen", "openai"},
		FallbackModels:      []string{"gpt-4", "gpt-3.5"},
		IsSystemDefault:     false,
	}

	if rule.GetUserId() != "user-123" {
		t.Errorf("UserId = %v, want user-123", rule.GetUserId())
	}
	if rule.GetIsSystemDefault() != false {
		t.Errorf("IsSystemDefault = %v, want false", rule.GetIsSystemDefault())
	}
	if len(rule.GetFallbackProviderIds()) != 2 {
		t.Errorf("FallbackProviderIds length = %v, want 2", len(rule.GetFallbackProviderIds()))
	}
}

func TestResolveRouteRequest_UserId(t *testing.T) {
	req := &ResolveRouteRequest{
		Model:             "ollama:llama2",
		AuthorizedModels:  []string{"ollama:llama2", "opencode_zen:gpt-4"},
		UserId:            "user-123",
	}

	if req.GetUserId() != "user-123" {
		t.Errorf("UserId = %v, want user-123", req.GetUserId())
	}
}

func TestRouteResult_FallbackProviderIds(t *testing.T) {
	result := &RouteResult{
		ProviderId:         "ollama",
		AdapterType:        "ollama",
		FallbackProviderIds: []string{"opencode_zen"},
		FallbackModels:      []string{"gpt-4"},
	}

	if len(result.GetFallbackProviderIds()) != 1 {
		t.Errorf("FallbackProviderIds length = %v, want 1", len(result.GetFallbackProviderIds()))
	}
}

func TestRoutingRule_Serialization(t *testing.T) {
	original := &RoutingRule{
		Id:                 "rule-1",
		UserId:              "user-123",
		ModelPattern:        "ollama:*",
		ProviderId:          "ollama",
		FallbackProviderIds: []string{"opencode_zen"},
		IsSystemDefault:     true,
	}

	data, err := proto.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	restored := &RoutingRule{}
	err = proto.Unmarshal(data, restored)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if restored.GetUserId() != original.GetUserId() {
		t.Errorf("UserId after round-trip = %v, want %v", restored.GetUserId(), original.GetUserId())
	}
	if restored.GetIsSystemDefault() != original.GetIsSystemDefault() {
		t.Errorf("IsSystemDefault after round-trip = %v, want %v", restored.GetIsSystemDefault(), original.GetIsSystemDefault())
	}
}
