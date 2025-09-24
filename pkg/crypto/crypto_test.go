package crypto

import (
	"bytes"
	"testing"
)

func TestSignAndVerify(t *testing.T) {
	// 1. Generate a key pair
	pubKey, privKey, err := GenerateKeys()
	if err != nil {
		t.Fatalf("Failed to generate keys: %v", err)
	}

	// 2. Create a message and sign it
	message := []byte("This is a test message for the DigiSocialBlock.")
	signature := Sign(privKey, message)

	// 3. Verify the signature with the correct public key and message
	if !Verify(pubKey, message, signature) {
		t.Fatal("Expected signature to be valid, but it was not")
	}
}

func TestVerify_InvalidMessage(t *testing.T) {
	// 1. Generate a key pair
	pubKey, privKey, err := GenerateKeys()
	if err != nil {
		t.Fatalf("Failed to generate keys: %v", err)
	}

	// 2. Create a message and sign it
	message := []byte("This is a test message.")
	signature := Sign(privKey, message)

	// 3. Attempt to verify with a different message
	invalidMessage := []byte("This is a different message.")
	if Verify(pubKey, invalidMessage, signature) {
		t.Fatal("Expected signature to be invalid for a different message, but it was valid")
	}
}

func TestVerify_InvalidPublicKey(t *testing.T) {
	// 1. Generate two key pairs
	_, privKey1, err := GenerateKeys()
	if err != nil {
		t.Fatalf("Failed to generate keys for signer: %v", err)
	}
	pubKey2, _, err := GenerateKeys()
	if err != nil {
		t.Fatalf("Failed to generate keys for verifier: %v", err)
	}

	// 2. Create a message and sign it with the first key
	message := []byte("This is a test message.")
	signature := Sign(privKey1, message)

	// 3. Attempt to verify with the second public key
	if Verify(pubKey2, message, signature) {
		t.Fatal("Expected signature to be invalid with a different public key, but it was valid")
	}
}

func TestVerify_CorruptedSignature(t *testing.T) {
	// 1. Generate a key pair
	pubKey, privKey, err := GenerateKeys()
	if err != nil {
		t.Fatalf("Failed to generate keys: %v", err)
	}

	// 2. Create a message and sign it
	message := []byte("This is a test message.")
	signature := Sign(privKey, message)

	// 3. Corrupt the signature
	corruptedSignature := bytes.Clone(signature)
	corruptedSignature[0] = ^corruptedSignature[0] // Flip the bits of the first byte

	// 4. Attempt to verify with the corrupted signature
	if Verify(pubKey, message, corruptedSignature) {
		t.Fatal("Expected corrupted signature to be invalid, but it was valid")
	}
}