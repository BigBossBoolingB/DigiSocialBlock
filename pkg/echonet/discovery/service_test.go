package discovery

import (
	"fmt"
	"testing"
	"time"

	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/crypto"
	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/echonet/peerstore"
	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/echonet/types"
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

func TestDiscovery_Announce(t *testing.T) {
	var receivedRequests []*types.AnnounceRequest
	mockBroadcaster := func(peerAddress string, request *types.AnnounceRequest) error {
		receivedRequests = append(receivedRequests, request)
		return nil
	}

	service := NewService(nil, nil, []string{"/bootstrap"}, mockBroadcaster, nil)
	_, privKey, _ := crypto.GenerateKeys()
	localPeer := types.PeerInfo{UserID: "local-peer-1"}

	err := service.Announce(localPeer, privKey)
	if err != nil {
		t.Fatalf("Announce failed: %v", err)
	}

	if len(receivedRequests) != 1 {
		t.Fatalf("Expected 1 request, got %d", len(receivedRequests))
	}
}

func TestDiscovery_ProcessAnnouncement(t *testing.T) {
	pubKey, privKey, _ := crypto.GenerateKeys()
	mockUser := &identity.NexusUserObjectV1{UserId: "announcer-1", PublicKey: pubKey}
	mockIS := &MockIdentityService{User: mockUser}
	ps := peerstore.New()
	service := NewService(mockIS, ps, nil, nil, nil)

	req := &types.AnnounceRequest{
		PeerInfo:  types.PeerInfo{UserID: "announcer-1"},
		Timestamp: time.Now(),
	}
	payload := []byte(fmt.Sprintf("%s|%d", req.PeerInfo.UserID, req.Timestamp.UnixNano()))
	req.Signature = crypto.Sign(privKey, payload)

	err := service.ProcessAnnouncement(req)
	if err != nil {
		t.Fatalf("ProcessAnnouncement failed: %v", err)
	}

	_, found := ps.Get("announcer-1")
	if !found {
		t.Fatal("Peer was not added to store")
	}
}

func TestDiscovery_DiscoverPeers(t *testing.T) {
	discoveredPeer := types.PeerInfo{UserID: "discovered-peer-1"}
	mockRequester := func(peerAddress string, request *types.FindPeersRequest) (*types.FindPeersResponse, error) {
		return &types.FindPeersResponse{
			Peers: []types.PeerInfo{discoveredPeer},
		}, nil
	}

	ps := peerstore.New()
	service := NewService(nil, ps, []string{"/bootstrap"}, nil, mockRequester)

	err := service.DiscoverPeers("#testing")
	if err != nil {
		t.Fatalf("DiscoverPeers failed: %v", err)
	}

	_, found := ps.Get("discovered-peer-1")
	if !found {
		t.Fatal("Peer was not added to store")
	}
}