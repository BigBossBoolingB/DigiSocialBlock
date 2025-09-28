package identity

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/crypto"
	"github.com/google/uuid"
)

// Service manages user identities.
// For the MVP, it uses a simple in-memory map for storage.
type Service struct {
	mu            sync.RWMutex
	usersByID     map[string]*NexusUserObjectV1
	usersByHandle map[string]*NexusUserObjectV1
}

// NewService creates and returns a new Identity Service.
func NewService() *Service {
	return &Service{
		usersByID:     make(map[string]*NexusUserObjectV1),
		usersByHandle: make(map[string]*NexusUserObjectV1),
	}
}

// CreateUser creates a new user, stores it, and returns the user object.
func (s *Service) CreateUser(publicKey []byte, handle string) (*NexusUserObjectV1, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.usersByHandle[handle]; exists {
		return nil, fmt.Errorf("handle '%s' is already taken", handle)
	}

	id := uuid.New().String()
	user := &NexusUserObjectV1{
		UserId:    id,
		PublicKey: publicKey,
		Handle:    handle,
		CreatedAt: time.Now(),
	}

	s.usersByID[user.UserId] = user
	s.usersByHandle[user.Handle] = user

	return user, nil
}

// GetUserByID retrieves a user by their unique ID.
func (s *Service) GetUserByID(id string) (*NexusUserObjectV1, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.usersByID[id]
	if !exists {
		return nil, fmt.Errorf("user with ID '%s' not found", id)
	}
	return user, nil
}

// GetUserByHandle retrieves a user by their unique handle.
func (s *Service) GetUserByHandle(handle string) (*NexusUserObjectV1, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.usersByHandle[handle]
	if !exists {
		return nil, fmt.Errorf("user with handle '%s' not found", handle)
	}
	return user, nil
}

// GenerateNonce creates a secure random string to be used as a challenge for authentication.
func (s *Service) GenerateNonce() (string, error) {
	nonceBytes := make([]byte, 32)
	_, err := rand.Read(nonceBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate random nonce: %w", err)
	}
	return hex.EncodeToString(nonceBytes), nil
}

// Authenticate verifies a user's signature over a given nonce.
func (s *Service) Authenticate(userID string, signature, nonce []byte) (bool, error) {
	user, err := s.GetUserByID(userID)
	if err != nil {
		return false, err
	}
	return crypto.Verify(user.PublicKey, nonce, signature), nil
}