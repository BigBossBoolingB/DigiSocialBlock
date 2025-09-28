package peerstore

import (
	"fmt"
	"sync"

	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/echonet/types"
)

// PeerStore is a thread-safe in-memory store for discovered peers.
type PeerStore struct {
	mu    sync.RWMutex
	peers map[string]*types.PeerInfo // Keyed by Peer's UserID
}

// New creates and returns a new PeerStore.
func New() *PeerStore {
	return &PeerStore{
		peers: make(map[string]*types.PeerInfo),
	}
}

// Add adds a new peer to the store or updates an existing one.
func (s *PeerStore) Add(peer *types.PeerInfo) error {
	if peer == nil || peer.UserID == "" {
		return fmt.Errorf("cannot add nil or invalid peer")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.peers[peer.UserID] = peer
	return nil
}

// Get retrieves a peer from the store by its user ID.
func (s *PeerStore) Get(userID string) (*types.PeerInfo, bool) {
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
func (s *PeerStore) List() []*types.PeerInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]*types.PeerInfo, 0, len(s.peers))
	for _, peer := range s.peers {
		list = append(list, peer)
	}
	return list
}