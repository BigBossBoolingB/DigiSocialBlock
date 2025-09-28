package echonet

import (
	"fmt"
	"net/rpc"
	"testing"
	"time"

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
	ds := discovery.NewService(mockIS, ps, nil, nil, nil)
	server, err := NewServer(ds)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// 2. Start the server
	if err := server.Start("127.0.0.1:0"); err != nil {
		t.Fatalf("Server failed to start: %v", err)
	}
	defer server.Stop()
	time.Sleep(100 * time.Millisecond) // Allow server to start

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
	ds := discovery.NewService(nil, ps, nil, nil, nil)
	server, err := NewServer(ds)
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