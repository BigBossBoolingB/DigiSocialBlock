package content

import (
	"fmt"
	"sync"

	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/crypto"
	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/identity"
	"github.com/google/uuid"
)

// IdentityService defines the interface for the functions we need from the identity package.
type IdentityService interface {
	GetUserByID(id string) (*identity.NexusUserObjectV1, error)
}

// Service manages content objects.
type Service struct {
	mu              sync.RWMutex
	identityService IdentityService
	contentByID     map[string]*ContentObjectV1
	contentByAuthor map[string][]*ContentObjectV1
}

// NewService creates a new Content Service.
func NewService(is IdentityService) *Service {
	return &Service{
		identityService: is,
		contentByID:     make(map[string]*ContentObjectV1),
		contentByAuthor: make(map[string][]*ContentObjectV1),
	}
}

// CreateContent validates a new content object and stores it.
// It ensures that the author's signature over the content hash is valid.
func (s *Service) CreateContent(req *ContentObjectV1) (*ContentObjectV1, error) {
	if req == nil {
		return nil, fmt.Errorf("content request cannot be nil")
	}

	// 1. Fetch the author's user object to get their public key.
	author, err := s.identityService.GetUserByID(req.AuthorUserID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify author identity: %w", err)
	}

	// 2. Verify the signature.
	if !crypto.Verify(author.PublicKey, req.ContentBodyHash, req.Signature) {
		return nil, fmt.Errorf("invalid signature for content")
	}

	// 3. If valid, finalize and store the content object.
	s.mu.Lock()
	defer s.mu.Unlock()

	req.ContentID = uuid.New().String() // Assign a new unique ID

	s.contentByID[req.ContentID] = req
	s.contentByAuthor[req.AuthorUserID] = append(s.contentByAuthor[req.AuthorUserID], req)

	return req, nil
}

// GetContentByID retrieves a content object by its unique ID.
func (s *Service) GetContentByID(id string) (*ContentObjectV1, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	content, found := s.contentByID[id]
	if !found {
		return nil, fmt.Errorf("content with ID '%s' not found", id)
	}
	return content, nil
}

// GetContentByAuthor retrieves all content objects created by a specific author.
func (s *Service) GetContentByAuthor(authorID string) ([]*ContentObjectV1, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	content, found := s.contentByAuthor[authorID]
	if !found {
		// Return an empty slice instead of an error if the author has no content
		return []*ContentObjectV1{}, nil
	}
	return content, nil
}