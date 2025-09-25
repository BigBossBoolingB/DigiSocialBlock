package routing

import (
	"fmt"

	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/echonet/peerstore"
	pb "github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/proto/echonet/v1"
)

// Service is responsible for determining routes to other peers in the network.
type Service struct {
	peerStore *peerstore.PeerStore
}

// NewService creates a new Routing Service.
func NewService(ps *peerstore.PeerStore) *Service {
	return &Service{
		peerStore: ps,
	}
}

// FindRoute looks up the connection information for a target peer.
// In this initial version, a "route" is simply the PeerInfo object itself,
// which contains the peer's multiaddress.
func (s *Service) FindRoute(targetUserID string) (*pb.PeerInfo, error) {
	if targetUserID == "" {
		return nil, fmt.Errorf("target user ID cannot be empty")
	}

	peer, found := s.peerStore.Get(targetUserID)
	if !found {
		return nil, fmt.Errorf("no route found for peer %s", targetUserID)
	}

	return peer, nil
}