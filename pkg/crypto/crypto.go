package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"fmt"
)

// GenerateKeys creates a new ed25519 public/private key pair.
func GenerateKeys() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate ed25519 key pair: %w", err)
	}
	return pubKey, privKey, nil
}

// Sign creates a cryptographic signature for a given message using a private key.
func Sign(privateKey ed25519.PrivateKey, message []byte) []byte {
	return ed25519.Sign(privateKey, message)
}

// Verify checks if a signature is valid for a given message and public key.
func Verify(publicKey ed25519.PublicKey, message []byte, signature []byte) bool {
	return ed25519.Verify(publicKey, message, signature)
}