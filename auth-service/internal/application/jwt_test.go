package application

import (
	"testing"
	"time"
)

func TestGenerateJWT(t *testing.T) {
	userID := "usr_test123"
	email := "test@example.com"
	role := "admin"

	token, err := GenerateJWT(userID, email, role, 24*time.Hour)
	if err != nil {
		t.Fatalf("GenerateJWT failed: %v", err)
	}
	if token == "" {
		t.Fatal("GenerateJWT returned empty token")
	}
}

func TestValidateJWT(t *testing.T) {
	userID := "usr_test123"
	email := "test@example.com"
	role := "admin"

	token, err := GenerateJWT(userID, email, role, 24*time.Hour)
	if err != nil {
		t.Fatalf("GenerateJWT failed: %v", err)
	}

	claims, err := ValidateJWT(token)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}
	if claims.Subject != userID {
		t.Errorf("ValidateJWT subject = %s, want %s", claims.Subject, userID)
	}
}

func TestValidateJWT_Invalid(t *testing.T) {
	_, err := ValidateJWT("invalid.token.string")
	if err == nil {
		t.Fatal("ValidateJWT should fail for invalid token")
	}
}

func TestValidateJWT_Expired(t *testing.T) {
	userID := "usr_test123"
	token, _ := GenerateJWT(userID, "test@example.com", "admin", -1*time.Hour)

	_, err := ValidateJWT(token)
	if err == nil {
		t.Fatal("ValidateJWT should fail for expired token")
	}
}