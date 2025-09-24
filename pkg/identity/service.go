package identity

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/crypto"
	pb "github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/proto"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Service manages user identities.
// For the MVP, it uses a simple in-memory map for storage.
type Service struct {
	mu            sync.RWMutex
	usersByID     map[string]*pb.NexusUserObjectV1
	usersByHandle map[string]*pb.NexusUserObjectV1
}

// NewService creates and returns a new Identity Service.
func NewService() *Service {
	return &Service{
		usersByID:     make(map[string]*pb.NexusUserObjectV1),
		usersByHandle: make(map[string]*pb.NexusUserObjectV1),
	}
}

// CreateUser creates a new user, stores it, and returns the user object.
func (s *Service) CreateUser(publicKey []byte, handle string) (*pb.NexusUserObjectV1, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Validate that the handle is not already taken
	if _, exists := s.usersByHandle[handle]; exists {
		return nil, fmt.Errorf("handle '%s' is already taken", handle)
	}

	// Create the new user object
	id := uuid.New().String()
	user := &pb.NexusUserObjectV1{
		UserId:    id,
		PublicKey: publicKey,
		Handle:    handle,
		CreatedAt: timestamppb.New(time.Now()),
	}

	// Store the user in our in-memory maps
	s.usersByID[user.UserId] = user
	s.usersByHandle[user.Handle] = user

	return user, nil
}

// GetUserByID retrieves a user by their unique ID.
func (s *Service) GetUserByID(id string) (*pb.NexusUserObjectV1, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.usersByID[id]
	if !exists {
		return nil, fmt.Errorf("user with ID '%s' not found", id)
	}
	return user, nil
}

// GenerateNonce creates a secure random string to be used as a challenge for authentication.
func (s *Service) GenerateNonce() (string, error) {
	// A 32-byte random nonce is sufficient for preventing replay attacks.
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
		return false, err // User not found
	}

	// The message that the user should have signed is the nonce itself.
	isValid := crypto.Verify(user.PublicKey, nonce, signature)
	return isValid, nil
}

// GetUserByHandle retrieves a user by their unique handle.
func (s *Service) GetUserByHandle(handle string) (*pb.NexusUserObjectV1, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, exists := s.usersByHandle[handle]
	if !exists {
		return nil, fmt.Errorf("user with handle '%s' not found", handle)
	}
	return user, nil
}