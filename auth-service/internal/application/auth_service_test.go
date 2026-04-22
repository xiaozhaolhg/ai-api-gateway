package application

import (
	"testing"
)

func TestAuthService_GenerateAPIKey(t *testing.T) {
	service := &AuthService{}

	// Generate multiple keys to ensure they are unique
	keys := make(map[string]bool)
	for i := 0; i < 100; i++ {
		key, err := service.GenerateAPIKey()
		if err != nil {
			t.Fatalf("GenerateAPIKey() error = %v", err)
		}
		if key == "" {
			t.Error("GenerateAPIKey() returned empty string")
		}
		if len(key) < 3 {
			t.Errorf("GenerateAPIKey() returned key shorter than expected: %d", len(key))
		}
		if key[:3] != "sk-" {
			t.Errorf("GenerateAPIKey() returned key with wrong prefix: %s", key[:3])
		}
		if keys[key] {
			t.Error("GenerateAPIKey() generated duplicate key")
		}
		keys[key] = true
	}
}

func TestAuthService_HashAPIKey(t *testing.T) {
	service := &AuthService{}

	tests := []struct {
		name    string
		apiKey  string
		wantErr bool
	}{
		{
			name:    "valid API key",
			apiKey:  "sk-test123456789",
			wantErr: false,
		},
		{
			name:    "empty API key",
			apiKey:  "",
			wantErr: false, // Hashing empty string is valid
		},
		{
			name:    "long API key",
			apiKey:  "sk-verylongkeywithlotsofrandomcharacters1234567890",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash := service.HashAPIKey(tt.apiKey)
			if hash == "" {
				t.Error("HashAPIKey() returned empty string")
			}
			if len(hash) != 64 { // SHA-256 produces 64 hex characters
				t.Errorf("HashAPIKey() returned hash of wrong length: %d", len(hash))
			}

			// Verify that the same key produces the same hash
			hash2 := service.HashAPIKey(tt.apiKey)
			if hash != hash2 {
				t.Error("HashAPIKey() is not deterministic")
			}

			// Verify that different keys produce different hashes
			if tt.apiKey != "" {
				differentKey := tt.apiKey + "different"
				hash3 := service.HashAPIKey(differentKey)
				if hash == hash3 {
					t.Error("HashAPIKey() produced same hash for different keys")
				}
			}
		})
	}
}

func TestAuthService_HashAPIKey_Consistency(t *testing.T) {
	service := &AuthService{}

	apiKey := "sk-test123456789"
	hash1 := service.HashAPIKey(apiKey)
	hash2 := service.HashAPIKey(apiKey)
	hash3 := service.HashAPIKey(apiKey)

	if hash1 != hash2 || hash2 != hash3 {
		t.Error("HashAPIKey() should produce consistent hashes for the same input")
	}
}
