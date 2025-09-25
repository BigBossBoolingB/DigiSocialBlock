package discovery

import (
	"testing"

	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/crypto"
	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/echonet/peerstore"
	pb "github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/proto/echonet/v1"
	"google.golang.org/protobuf/proto"
)

func TestDiscovery_Announce(t *testing.T) {
	// 1. Setup mock broadcaster and state
	var receivedRequests []*pb.AnnounceRequest
	mockBroadcaster := func(peerAddress string, request *pb.AnnounceRequest) error {
		receivedRequests = append(receivedRequests, request)
		return nil
	}

	bootstrapPeers := []string{"/ip4/127.0.0.1/tcp/4001", "/ip4/127.0.0.1/tcp/4002"}
	service := NewService(nil, bootstrapPeers, mockBroadcaster, nil)

	// 2. Create local peer info and keys
	pubKey, privKey, err := crypto.GenerateKeys()
	if err != nil {
		t.Fatalf("Failed to generate keys: %v", err)
	}
	localPeer := &pb.PeerInfo{
		UserId:       "local-peer-1",
		Multiaddress: "/ip4/127.0.0.1/tcp/5001",
		Frequencies:  []string{"#testing"},
	}

	// 3. Call the Announce method
	err = service.Announce(localPeer, privKey)
	if err != nil {
		t.Fatalf("Announce failed unexpectedly: %v", err)
	}

	// 4. Verify the results
	if len(receivedRequests) != len(bootstrapPeers) {
		t.Fatalf("Expected broadcaster to be called %d times, but got %d", len(bootstrapPeers), len(receivedRequests))
	}

	// Check the content of the first received request
	req := receivedRequests[0]
	if req.PeerInfo.UserId != localPeer.UserId {
		t.Errorf("Expected announced user ID to be %s, but got %s", localPeer.UserId, req.PeerInfo.UserId)
	}

	// Verify the signature
	// Re-create the signed payload to verify against
	payload, err := proto.Marshal(&pb.AnnounceRequest{
		PeerInfo:  req.PeerInfo,
		Timestamp: req.Timestamp,
	})
	if err != nil {
		t.Fatalf("Failed to marshal payload for verification: %v", err)
	}

	if !crypto.Verify(pubKey, payload, req.Signature) {
		t.Error("Signature on announce request is invalid")
	}
}

func TestDiscovery_Announce_InvalidPeer(t *testing.T) {
	service := NewService(nil, nil, nil, nil)
	_, privKey, _ := crypto.GenerateKeys()

	err := service.Announce(nil, privKey)
	if err == nil {
		t.Error("Expected error when announcing with nil peer, but got nil")
	}

	err = service.Announce(&pb.PeerInfo{UserId: ""}, privKey)
	if err == nil {
		t.Error("Expected error when announcing with empty peer ID, but got nil")
	}
}

func TestDiscovery_DiscoverPeers(t *testing.T) {
	// 1. Setup mock requester and state
	discoveredPeer := &pb.PeerInfo{UserId: "discovered-peer-1", Multiaddress: "/ip4/10.0.0.1/tcp/4001"}
	mockRequester := func(peerAddress string, request *pb.FindPeersRequest) (*pb.FindPeersResponse, error) {
		// Simulate a successful response from a bootstrap peer
		return &pb.FindPeersResponse{
			Peers: []*pb.PeerInfo{discoveredPeer},
		}, nil
	}

	ps := peerstore.New()
	bootstrapPeers := []string{"/ip4/127.0.0.1/tcp/4001"}
	service := NewService(ps, bootstrapPeers, nil, mockRequester)

	// 2. Call the DiscoverPeers method
	err := service.DiscoverPeers("#testing")
	if err != nil {
		t.Fatalf("DiscoverPeers failed unexpectedly: %v", err)
	}

	// 3. Verify that the peer store was populated
	retrievedPeer, found := ps.Get("discovered-peer-1")
	if !found {
		t.Fatal("Expected to find discovered peer in peer store, but it was not found")
	}
	if retrievedPeer.Multiaddress != discoveredPeer.Multiaddress {
		t.Errorf("Expected peer address to be %s, but got %s", discoveredPeer.Multiaddress, retrievedPeer.Multiaddress)
	}
}