package application

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	pass := "testpassword123"
	hash, err := HashPassword(pass)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if hash == "" {
		t.Fatal("HashPassword returned empty string")
	}
	if hash == pass {
		t.Fatal("HashPassword did not hash the password")
	}
	if len(hash) < 50 {
		t.Fatalf("HashPassword hash too short: %d", len(hash))
	}
}

func TestCheckPassword(t *testing.T) {
	pass := "testpassword123"
	hash, _ := HashPassword(pass)

	if !CheckPassword(hash, pass) {
		t.Fatal("CheckPassword returned false for correct password")
	}
	if CheckPassword(hash, "wrongpassword") {
		t.Fatal("CheckPassword returned true for incorrect password")
	}
}

func TestCheckPassword_SameInput_DifferentHash(t *testing.T) {
	pass := "testpassword123"
	hash1, _ := HashPassword(pass)
	hash2, _ := HashPassword(pass)

	if hash1 == hash2 {
		t.Fatal("Same input produced same hash (salt not working)")
	}

	if !CheckPassword(hash1, pass) {
		t.Fatal("CheckPassword failed for hash1")
	}
	if !CheckPassword(hash2, pass) {
		t.Fatal("CheckPassword failed for hash2")
	}
}