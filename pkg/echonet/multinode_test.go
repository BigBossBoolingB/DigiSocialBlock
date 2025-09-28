package echonet

import (
	"fmt"
	"testing"
	"time"

	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/content"
	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/crypto"
	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/echonet/discovery"
	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/echonet/peerstore"
	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/echonet/types"
	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/identity"
)

// setupNode creates a full, functional node stack for testing.
// It takes an existing identity service to simulate a shared identity layer.
func setupNode(t *testing.T, idService *identity.Service) (*peerstore.PeerStore, *Server) {
	contentService := content.NewService(idService)
	ps := peerstore.New()

	// The discovery service is wired with the live network functions
	ds := discovery.NewService(idService, ps, []string{}, discovery.LiveAnnounceBroadcaster, discovery.LiveFindPeersRequester)

	server, err := NewServer(ds, contentService)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	if err := server.Start("127.0.0.1:0"); err != nil {
		t.Fatalf("Server failed to start: %v", err)
	}
	time.Sleep(100 * time.Millisecond) // Allow server to start

	return ps, server
}

func TestMultiNode_Announce(t *testing.T) {
	// 1. Create a single, shared identity service for all nodes
	idService := identity.NewService()

	// 2. Create two independent nodes that share the identity service
	_, server1 := setupNode(t, idService)
	defer server1.Stop()

	ps2, server2 := setupNode(t, idService)
	defer server2.Stop()

	// Node 1 needs to know about Node 2 to announce to it.
	// We re-create its discovery service with the bootstrap peer info.
	server1.API.Discovery = discovery.NewService(idService, peerstore.New(), []string{server2.Address()}, discovery.LiveAnnounceBroadcaster, nil)

	// 3. Create an identity on the shared identity service
	pubKey1, privKey1, _ := crypto.GenerateKeys()
	user1, err := idService.CreateUser(pubKey1, "node1")
	if err != nil {
		t.Fatalf("Failed to create user on node 1: %v", err)
	}

	// 3. Node 1 announces itself to Node 2
	peerInfo1 := types.PeerInfo{
		UserID:       user1.UserId,
		Multiaddress: server1.Address(),
		Frequencies:  []string{"#multinode-test"},
	}
	err = server1.API.Discovery.Announce(peerInfo1, privKey1)
	if err != nil {
		t.Fatalf("Node 1 failed to announce to Node 2: %v", err)
	}

	// 4. Verify that Node 2 now has Node 1 in its peer store
	time.Sleep(100 * time.Millisecond) // Allow for processing

	retrievedPeer, found := ps2.Get(user1.UserId)
	if !found {
		t.Fatal("Node 2 did not have Node 1 in its peer store after announcement")
	}

	if retrievedPeer.Multiaddress != server1.Address() {
		t.Errorf("Expected peer address to be %s, but got %s", server1.Address(), retrievedPeer.Multiaddress)
	}
	fmt.Printf("Success! Node 2 discovered Node 1 at address %s\n", retrievedPeer.Multiaddress)
}