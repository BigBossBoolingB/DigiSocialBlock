package routing

import (
	"testing"

	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/echonet/peerstore"
	pb "github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/proto/echonet/v1"
)

func TestRouting_FindRoute_Success(t *testing.T) {
	// 1. Setup
	ps := peerstore.New()
	service := NewService(ps)

	// 2. Pre-populate the peer store
	testPeer := &pb.PeerInfo{UserId: "target-peer", Multiaddress: "/ip4/127.0.0.1/tcp/4001"}
	if err := ps.Add(testPeer); err != nil {
		t.Fatalf("Failed to add peer to peerstore for testing: %v", err)
	}

	// 3. Call FindRoute for the existing peer
	route, err := service.FindRoute("target-peer")
	if err != nil {
		t.Fatalf("FindRoute failed unexpectedly: %v", err)
	}

	// 4. Verify the result
	if route.UserId != testPeer.UserId {
		t.Errorf("Expected route to have user ID %s, but got %s", testPeer.UserId, route.UserId)
	}
	if route.Multiaddress != testPeer.Multiaddress {
		t.Errorf("Expected route to have address %s, but got %s", testPeer.Multiaddress, route.Multiaddress)
	}
}

func TestRouting_FindRoute_NotFound(t *testing.T) {
	// 1. Setup
	ps := peerstore.New()
	service := NewService(ps)

	// 2. Call FindRoute for a non-existent peer
	_, err := service.FindRoute("non-existent-peer")
	if err == nil {
		t.Fatal("Expected an error when finding route to non-existent peer, but got nil")
	}
}

func TestRouting_FindRoute_EmptyID(t *testing.T) {
	// 1. Setup
	ps := peerstore.New()
	service := NewService(ps)

	// 2. Call FindRoute with an empty user ID
	_, err := service.FindRoute("")
	if err == nil {
		t.Fatal("Expected an error when finding route with empty user ID, but got nil")
	}
}