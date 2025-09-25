package echonet

import (
	"testing"

	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/echonet/discovery"
	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/echonet/peerstore"
	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/echonet/routing"
	pb "github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/proto/echonet/v1"
)

// TestIntegration_DiscoveryToRouting verifies the end-to-end flow of
// discovering a peer and then finding a route to that peer.
func TestIntegration_DiscoveryToRouting(t *testing.T) {
	// 1. Setup all the components
	ps := peerstore.New()
	rs := routing.NewService(ps)

	// The peer we expect to discover
	discoveredPeer := &pb.PeerInfo{
		UserId:       "discovered-peer-1",
		Multiaddress: "/ip4/10.0.0.1/tcp/4001",
		Frequencies:  []string{"#integration-test"},
	}

	// Mock the network requester for the discovery service
	mockRequester := func(peerAddress string, request *pb.FindPeersRequest) (*pb.FindPeersResponse, error) {
		// When the discovery service asks for peers, return our test peer
		return &pb.FindPeersResponse{
			Peers: []*pb.PeerInfo{discoveredPeer},
		}, nil
	}

	bootstrapPeers := []string{"/ip4/127.0.0.1/tcp/4001"}
	ds := discovery.NewService(ps, bootstrapPeers, nil, mockRequester)

	// 2. Execute the discovery process
	err := ds.DiscoverPeers("#integration-test")
	if err != nil {
		t.Fatalf("Discovery process failed unexpectedly: %v", err)
	}

	// 3. Execute the routing process for the discovered peer
	route, err := rs.FindRoute("discovered-peer-1")
	if err != nil {
		t.Fatalf("Routing process failed unexpectedly: %v", err)
	}

	// 4. Verify the final outcome
	if route == nil {
		t.Fatal("Expected to find a route, but got nil")
	}
	if route.UserId != discoveredPeer.UserId {
		t.Errorf("Expected route to have user ID %s, but got %s", discoveredPeer.UserId, route.UserId)
	}
	if route.Multiaddress != discoveredPeer.Multiaddress {
		t.Errorf("Expected route to have address %s, but got %s", discoveredPeer.Multiaddress, route.Multiaddress)
	}

	// Also verify that the peer is actually in the peer store
	_, found := ps.Get("discovered-peer-1")
	if !found {
		t.Fatal("Peer was not found in the peer store after discovery")
	}
}