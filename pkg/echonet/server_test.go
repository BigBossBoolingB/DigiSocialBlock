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

func TestServer_RPC_Integration(t *testing.T) {
	// 1. Setup keys and a mock identity service for the server to use
	pubKey, privKey, _ := crypto.GenerateKeys()
	mockUser := &identity.NexusUserObjectV1{UserId: "live-client-1", PublicKey: pubKey}
	mockIS := &MockIdentityService{User: mockUser}

	// 2. Setup the full dependency chain for the server
	ps := peerstore.New()
	ds := discovery.NewService(mockIS, ps, nil, nil, nil)
	server, err := NewServer(ds)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// 3. Start the RPC server on a random available port
	go server.Start("127.0.0.1:0") // Use port 0 to get a random free port
	defer server.Stop()
	time.Sleep(100 * time.Millisecond) // Allow server to start
	serverAddr := server.Address()
	if serverAddr == "" {
		t.Fatal("Server did not start and get an address")
	}

	// 4. Create an RPC client and connect to the server
	client, err := rpc.Dial("tcp", serverAddr)
	if err != nil {
		t.Fatalf("Failed to dial RPC server: %v", err)
	}
	defer client.Close()

	// 5. Prepare the arguments for the RPC call
	args := &types.AnnounceRequest{
		PeerInfo: types.PeerInfo{
			UserID:       "live-client-1",
			Multiaddress: "/ip4/127.0.0.1/tcp/9001",
		},
		Timestamp: time.Now(),
	}
	// Sign the payload with the private key that corresponds to the public key in the mock identity service
	payload := []byte(fmt.Sprintf("%s|%d", args.PeerInfo.UserID, args.Timestamp.UnixNano()))
	args.Signature = crypto.Sign(privKey, payload)

	var reply types.AnnounceResponse

	// 6. Make the RPC call
	err = client.Call("EchoNetAPI.Announce", args, &reply)
	if err != nil {
		t.Fatalf("RPC call failed: %v", err)
	}

	// 7. Verify the response and side-effects
	expectedMsg := "Announcement for live-client-1 processed."
	if reply.ConfirmationMessage != expectedMsg {
		t.Errorf("Expected reply '%s', but got '%s'", expectedMsg, reply.ConfirmationMessage)
	}

	_, found := ps.Get("live-client-1")
	if !found {
		t.Fatal("Expected peer to be added to the peer store, but it was not")
	}
}