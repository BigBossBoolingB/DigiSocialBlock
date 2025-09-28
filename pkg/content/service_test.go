package content

import (
	"crypto/sha256"
	"fmt"
	"testing"

	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/crypto"
	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/identity"
)

// MockIdentityService is a mock implementation of the IdentityService interface.
type MockIdentityService struct {
	User *identity.NexusUserObjectV1
	Err  error
}

// GetUserByID implements the IdentityService interface for the mock.
func (m *MockIdentityService) GetUserByID(id string) (*identity.NexusUserObjectV1, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	if m.User != nil && m.User.UserId == id {
		return m.User, nil
	}
	return nil, fmt.Errorf("user not found")
}

func TestContent_CreateContent_Success(t *testing.T) {
	// 1. Setup
	pubKey, privKey, _ := crypto.GenerateKeys()
	mockUser := &identity.NexusUserObjectV1{UserId: "author-1", PublicKey: pubKey}
	mockIS := &MockIdentityService{User: mockUser}
	service := NewService(mockIS)

	// 2. Create content and sign its hash
	contentBody := []byte("This is the core content.")
	contentHash := sha256.Sum256(contentBody)
	signature := crypto.Sign(privKey, contentHash[:])

	req := &ContentObjectV1{
		AuthorUserID:    "author-1",
		ContentBodyHash: contentHash[:],
		Signature:       signature,
	}

	// 3. Call CreateContent
	content, err := service.CreateContent(req)
	if err != nil {
		t.Fatalf("CreateContent failed unexpectedly: %v", err)
	}

	// 4. Verify results
	if content.ContentID == "" {
		t.Error("Expected ContentID to be set, but it was empty")
	}
}

func TestContent_CreateContent_InvalidSignature(t *testing.T) {
	// 1. Setup
	pubKey, _, _ := crypto.GenerateKeys() // Real user key
	_, otherPrivKey, _ := crypto.GenerateKeys()   // Key to create invalid signature
	mockUser := &identity.NexusUserObjectV1{UserId: "author-1", PublicKey: pubKey}
	mockIS := &MockIdentityService{User: mockUser}
	service := NewService(mockIS)

	// 2. Create content and sign with the WRONG key
	contentBody := []byte("This is the core content.")
	contentHash := sha256.Sum256(contentBody)
	signature := crypto.Sign(otherPrivKey, contentHash[:])

	req := &ContentObjectV1{
		AuthorUserID:    "author-1",
		ContentBodyHash: contentHash[:],
		Signature:       signature,
	}

	// 3. Call CreateContent
	_, err := service.CreateContent(req)
	if err == nil {
		t.Fatal("Expected an error for invalid signature, but got nil")
	}
}

func TestContent_Getters(t *testing.T) {
	// 1. Setup
	pubKey, privKey, _ := crypto.GenerateKeys()
	mockUser := &identity.NexusUserObjectV1{UserId: "author-1", PublicKey: pubKey}
	mockIS := &MockIdentityService{User: mockUser}
	service := NewService(mockIS)

	// 2. Create some content
	hash := sha256.Sum256([]byte("content1"))
	sig := crypto.Sign(privKey, hash[:])
	req1, _ := service.CreateContent(&ContentObjectV1{AuthorUserID: "author-1", ContentBodyHash: hash[:], Signature: sig})

	hash = sha256.Sum256([]byte("content2"))
	sig = crypto.Sign(privKey, hash[:])
	_, _ = service.CreateContent(&ContentObjectV1{AuthorUserID: "author-1", ContentBodyHash: hash[:], Signature: sig})

	// 3. Test GetContentByID
	retrieved, err := service.GetContentByID(req1.ContentID)
	if err != nil {
		t.Fatalf("GetContentByID failed: %v", err)
	}
	if retrieved.ContentID != req1.ContentID {
		t.Error("GetContentByID returned wrong content")
	}

	// 4. Test GetContentByAuthor
	allContent, err := service.GetContentByAuthor("author-1")
	if err != nil {
		t.Fatalf("GetContentByAuthor failed: %v", err)
	}
	if len(allContent) != 2 {
		t.Errorf("Expected 2 content objects for author, but got %d", len(allContent))
	}
}