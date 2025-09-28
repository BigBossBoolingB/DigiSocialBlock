package identity

import (
	"encoding/hex"
	"testing"

	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/crypto"
)

func TestCreateUser(t *testing.T) {
	service := NewService()
	handle := "TheArchitect"
	publicKey, _, err := crypto.GenerateKeys()
	if err != nil {
		t.Fatalf("Failed to generate keys: %v", err)
	}

	user, err := service.CreateUser(publicKey, handle)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if user.Handle != handle {
		t.Errorf("Expected handle to be %s, but got %s", handle, user.Handle)
	}

	if user.UserId == "" {
		t.Error("Expected user ID to be non-empty")
	}

	if user.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
}

func TestCreateUser_DuplicateHandle(t *testing.T) {
	service := NewService()
	handle := "TheArchitect"
	publicKey1, _, _ := crypto.GenerateKeys()
	publicKey2, _, _ := crypto.GenerateKeys()

	_, err := service.CreateUser(publicKey1, handle)
	if err != nil {
		t.Fatalf("First user creation failed unexpectedly: %v", err)
	}

	_, err = service.CreateUser(publicKey2, handle)
	if err == nil {
		t.Fatal("Expected error when creating user with duplicate handle, but got nil")
	}
}

func TestGetUser(t *testing.T) {
	service := NewService()
	handle := "TheArchitect"
	publicKey, _, _ := crypto.GenerateKeys()

	createdUser, err := service.CreateUser(publicKey, handle)
	if err != nil {
		t.Fatalf("Failed to create user for testing: %v", err)
	}

	// Test GetUserByID
	userByID, err := service.GetUserByID(createdUser.UserId)
	if err != nil {
		t.Fatalf("Failed to get user by ID: %v", err)
	}
	if userByID.UserId != createdUser.UserId {
		t.Errorf("GetUserByID returned user with wrong ID")
	}

	// Test GetUserByHandle
	userByHandle, err := service.GetUserByHandle(handle)
	if err != nil {
		t.Fatalf("Failed to get user by handle: %v", err)
	}
	if userByHandle.Handle != handle {
		t.Errorf("GetUserByHandle returned user with wrong handle")
	}
}

func TestGetUser_NotFound(t *testing.T) {
	service := NewService()

	_, err := service.GetUserByID("non-existent-id")
	if err == nil {
		t.Fatal("Expected error when getting non-existent user by ID, but got nil")
	}

	_, err = service.GetUserByHandle("non-existent-handle")
	if err == nil {
		t.Fatal("Expected error when getting non-existent user by handle, but got nil")
	}
}

func TestAuthentication(t *testing.T) {
	service := NewService()

	// 1. Create a user and generate keys
	pubKey, privKey, err := crypto.GenerateKeys()
	if err != nil {
		t.Fatalf("Failed to generate keys: %v", err)
	}
	user, err := service.CreateUser(pubKey, "auth-user")
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// 2. Generate a nonce (challenge)
	nonceStr, err := service.GenerateNonce()
	if err != nil {
		t.Fatalf("Failed to generate nonce: %v", err)
	}
	nonceBytes, _ := hex.DecodeString(nonceStr)

	// 3. Sign the nonce with the user's private key
	signature := crypto.Sign(privKey, nonceBytes)

	// 4. Authenticate with the correct signature
	valid, err := service.Authenticate(user.UserId, signature, nonceBytes)
	if err != nil {
		t.Fatalf("Authentication failed unexpectedly: %v", err)
	}
	if !valid {
		t.Fatal("Expected authentication to be successful, but it failed")
	}

	// 5. Attempt to authenticate with a wrong nonce
	invalidNonce := []byte("this-is-a-wrong-nonce")
	invalid, err := service.Authenticate(user.UserId, signature, invalidNonce)
	if err != nil {
		t.Fatalf("Authentication with wrong nonce failed unexpectedly: %v", err)
	}
	if invalid {
		t.Fatal("Expected authentication with wrong nonce to fail, but it succeeded")
	}
}