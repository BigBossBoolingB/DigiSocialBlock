package echonet

import (
	"crypto/sha256"
	"fmt"
	"net/rpc"
	"testing"
	"time"

	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/content"
	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/crypto"
	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/echonet/discovery"
	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/echonet/peerstore"
	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/echonet/types"
	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/identity"
)

// MockIdentityService is a mock implementation of the discovery.IdentityService interface.
type MockIdentityService struct {
	User *identity.NexusUserObjectV1
	Err  error
}

// GetUserByID implements the required interface method for the mock.
func (m *MockIdentityService) GetUserByID(id string) (*identity.NexusUserObjectV1, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	if m.User != nil && m.User.UserId == id {
		return m.User, nil
	}
	return nil, fmt.Errorf("user not found")
}

func TestServer_RPC_Announce(t *testing.T) {
	// 1. Setup dependencies
	pubKey, privKey, _ := crypto.GenerateKeys()
	mockUser := &identity.NexusUserObjectV1{UserId: "live-client-1", PublicKey: pubKey}
	mockIS := &MockIdentityService{User: mockUser}
	ps := peerstore.New()
	cs := content.NewService(mockIS)
	ds := discovery.NewService(mockIS, ps, nil, nil, nil)
	server, err := NewServer(ds, cs)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// 2. Start the server
	if err := server.Start("127.0.0.1:0"); err != nil {
		t.Fatalf("Server failed to start: %v", err)
	}
	defer server.Stop()
	time.Sleep(100 * time.Millisecond)

	// 3. Create client and connect
	client, err := rpc.DialHTTP("tcp", server.Address())
	if err != nil {
		t.Fatalf("Failed to dial RPC server: %v", err)
	}
	defer client.Close()

	// 4. Prepare and make the RPC call
	args := &types.AnnounceRequest{
		PeerInfo: types.PeerInfo{UserID: "live-client-1"},
		Timestamp: time.Now(),
	}
	payload := []byte(fmt.Sprintf("%s|%d", args.PeerInfo.UserID, args.Timestamp.UnixNano()))
	args.Signature = crypto.Sign(privKey, payload)
	var reply types.AnnounceResponse
	err = client.Call("EchoNetAPI.Announce", args, &reply)
	if err != nil {
		t.Fatalf("RPC call failed: %v", err)
	}

	// 5. Verify results
	if _, found := ps.Get("live-client-1"); !found {
		t.Fatal("Expected peer to be in peer store after successful announcement")
	}
}

func TestServer_RPC_FindPeers(t *testing.T) {
	// 1. Setup dependencies
	ps := peerstore.New()
	cs := content.NewService(nil)
	ds := discovery.NewService(nil, ps, nil, nil, nil)
	server, err := NewServer(ds, cs)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}
	ps.Add(&types.PeerInfo{UserID: "peer1", Frequencies: []string{"#testing"}})
	ps.Add(&types.PeerInfo{UserID: "peer2", Frequencies: []string{"#testing"}})

	// 2. Start server
	if err := server.Start("127.0.0.1:0"); err != nil {
		t.Fatalf("Server failed to start: %v", err)
	}
	defer server.Stop()
	time.Sleep(100 * time.Millisecond)

	// 3. Create client and connect
	client, err := rpc.DialHTTP("tcp", server.Address())
	if err != nil {
		t.Fatalf("Failed to dial RPC server: %v", err)
	}
	defer client.Close()

	// 4. Make RPC call
	args := &types.FindPeersRequest{Frequency: "#testing"}
	var reply types.FindPeersResponse
	err = client.Call("EchoNetAPI.FindPeers", args, &reply)
	if err != nil {
		t.Fatalf("RPC call to FindPeers failed: %v", err)
	}

	// 5. Verify results
	if len(reply.Peers) != 2 {
		t.Fatalf("Expected to find 2 peers, but got %d", len(reply.Peers))
	}
}

func TestServer_RPC_PublishContent(t *testing.T) {
	// 1. Setup dependencies
	pubKey, privKey, _ := crypto.GenerateKeys()
	mockUser := &identity.NexusUserObjectV1{UserId: "author-1", PublicKey: pubKey}
	mockIS := &MockIdentityService{User: mockUser}
	ps := peerstore.New()
	cs := content.NewService(mockIS)
	ds := discovery.NewService(mockIS, ps, nil, nil, nil)
	server, err := NewServer(ds, cs)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// 2. Start server
	if err := server.Start("127.0.0.1:0"); err != nil {
		t.Fatalf("Server failed to start: %v", err)
	}
	defer server.Stop()
	time.Sleep(100 * time.Millisecond)

	// 3. Create client and connect
	client, err := rpc.DialHTTP("tcp", server.Address())
	if err != nil {
		t.Fatalf("Failed to dial RPC server: %v", err)
	}
	defer client.Close()

	// 4. Prepare the content and request
	contentBody := []byte("This is a test content.")
	contentHash := sha256.Sum256(contentBody)
	signature := crypto.Sign(privKey, contentHash[:])

	args := &types.PublishContentRequest{
		AuthorUserID:    "author-1",
		ContentBodyURI:  "mem://content-1",
		ContentBodyHash: contentHash[:],
		Signature:       signature,
	}
	var reply types.PublishContentResponse

	// 5. Make the RPC call
	err = client.Call("EchoNetAPI.PublishContent", args, &reply)
	if err != nil {
		t.Fatalf("RPC call to PublishContent failed: %v", err)
	}

	// 6. Verify the response and side-effects
	if reply.ContentID == "" {
		t.Error("Expected a ContentID in the response, but it was empty")
	}
	if _, err := cs.GetContentByID(reply.ContentID); err != nil {
		t.Error("Expected content to be in the content service, but it was not found")
	}
}