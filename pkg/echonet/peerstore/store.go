package peerstore

import (
	"fmt"
	"sync"

	pb "github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/proto/echonet/v1"
)

// PeerStore is a thread-safe in-memory store for discovered peers.
// It acts as a node's local routing table or address book.
type PeerStore struct {
	mu    sync.RWMutex
	peers map[string]*pb.PeerInfo // Keyed by Peer's user_id
}

// New creates and returns a new PeerStore.
func New() *PeerStore {
	return &PeerStore{
		peers: make(map[string]*pb.PeerInfo),
	}
}

// Add adds a new peer to the store or updates an existing one.
// It is safe for concurrent use.
func (s *PeerStore) Add(peer *pb.PeerInfo) error {
	if peer == nil || peer.UserId == "" {
		return fmt.Errorf("cannot add nil or invalid peer")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.peers[peer.UserId] = peer
	return nil
}

// Get retrieves a peer from the store by its user ID.
// It returns the PeerInfo and a boolean indicating if the peer was found.
func (s *PeerStore) Get(userID string) (*pb.PeerInfo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	peer, found := s.peers[userID]
	return peer, found
}

// Remove deletes a peer from the store by its user ID.
func (s *PeerStore) Remove(userID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.peers, userID)
}

// List returns a slice of all peers currently in the store.
func (s *PeerStore) List() []*pb.PeerInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]*pb.PeerInfo, 0, len(s.peers))
	for _, peer := range s.peers {
		list = append(list, peer)
	}
	return list
}