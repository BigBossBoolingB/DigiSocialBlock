package peerstore

import (
	"fmt"
	"sync"
	"testing"

	pb "github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/proto/echonet/v1"
)

func TestPeerStore_AddAndGet(t *testing.T) {
	store := New()
	peer := &pb.PeerInfo{UserId: "peer1", Multiaddress: "/ip4/127.0.0.1/tcp/4001"}

	err := store.Add(peer)
	if err != nil {
		t.Fatalf("Failed to add peer: %v", err)
	}

	retrievedPeer, found := store.Get("peer1")
	if !found {
		t.Fatal("Expected to find peer, but it was not found")
	}
	if retrievedPeer.Multiaddress != peer.Multiaddress {
		t.Errorf("Expected address %s, got %s", peer.Multiaddress, retrievedPeer.Multiaddress)
	}
}

func TestPeerStore_AddInvalid(t *testing.T) {
	store := New()
	if err := store.Add(nil); err == nil {
		t.Error("Expected error when adding nil peer, but got nil")
	}
	if err := store.Add(&pb.PeerInfo{UserId: ""}); err == nil {
		t.Error("Expected error when adding peer with empty ID, but got nil")
	}
}

func TestPeerStore_Update(t *testing.T) {
	store := New()
	peerV1 := &pb.PeerInfo{UserId: "peer1", Multiaddress: "/ip4/127.0.0.1/tcp/4001"}
	store.Add(peerV1)

	peerV2 := &pb.PeerInfo{UserId: "peer1", Multiaddress: "/ip4/127.0.0.1/tcp/4002"}
	store.Add(peerV2)

	retrievedPeer, _ := store.Get("peer1")
	if retrievedPeer.Multiaddress != peerV2.Multiaddress {
		t.Errorf("Expected address to be updated to %s, but got %s", peerV2.Multiaddress, retrievedPeer.Multiaddress)
	}
}

func TestPeerStore_GetNotFound(t *testing.T) {
	store := New()
	_, found := store.Get("non-existent-peer")
	if found {
		t.Fatal("Expected not to find peer, but it was found")
	}
}

func TestPeerStore_Remove(t *testing.T) {
	store := New()
	peer := &pb.PeerInfo{UserId: "peer1"}
	store.Add(peer)

	store.Remove("peer1")

	_, found := store.Get("peer1")
	if found {
		t.Fatal("Expected peer to be removed, but it was still found")
	}
}

func TestPeerStore_List(t *testing.T) {
	store := New()
	store.Add(&pb.PeerInfo{UserId: "peer1"})
	store.Add(&pb.PeerInfo{UserId: "peer2"})

	list := store.List()
	if len(list) != 2 {
		t.Errorf("Expected list of length 2, but got %d", len(list))
	}
}

func TestPeerStore_Concurrency(t *testing.T) {
	store := New()
	var wg sync.WaitGroup
	peerCount := 100

	// Concurrently add peers
	for i := 0; i < peerCount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			peerID := fmt.Sprintf("peer%d", i)
			peer := &pb.PeerInfo{UserId: peerID}
			store.Add(peer)
		}(i)
	}

	wg.Wait()

	// Check if all peers were added
	if len(store.List()) != peerCount {
		t.Errorf("Expected peer count to be %d, but got %d", peerCount, len(store.List()))
	}
}