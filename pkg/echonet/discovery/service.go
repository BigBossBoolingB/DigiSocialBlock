package discovery

import (
	"crypto/ed25519"
	"fmt"
	"time"

	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/crypto"
	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/echonet/peerstore"
	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/echonet/types"
	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/identity"
)

// Broadcaster defines the interface for sending an announcement to a specific peer.
type Broadcaster func(peerAddress string, request *types.AnnounceRequest) error

// Requester defines the interface for sending a request to a peer and getting a response.
type Requester func(peerAddress string, request *types.FindPeersRequest) (*types.FindPeersResponse, error)

// IdentityService defines the interface for the functions we need from the identity package.
type IdentityService interface {
	GetUserByID(id string) (*identity.NexusUserObjectV1, error)
}

// Service handles the peer discovery process.
type Service struct {
	identityService IdentityService
	peerStore       *peerstore.PeerStore
	bootstrapPeers  []string
	broadcaster     Broadcaster
	requester       Requester
}

// NewService creates a new discovery service.
func NewService(is IdentityService, ps *peerstore.PeerStore, bootstrapPeers []string, broadcaster Broadcaster, requester Requester) *Service {
	return &Service{
		identityService: is,
		peerStore:       ps,
		bootstrapPeers:  bootstrapPeers,
		broadcaster:     broadcaster,
		requester:       requester,
	}
}

// Announce broadcasts the presence of the local node to the network.
func (s *Service) Announce(localPeer types.PeerInfo, privateKey ed25519.PrivateKey) error {
	if localPeer.UserID == "" {
		return fmt.Errorf("cannot announce with invalid peer info")
	}

	request := &types.AnnounceRequest{
		PeerInfo:  localPeer,
		Timestamp: time.Now(),
	}

	payload := []byte(fmt.Sprintf("%s|%d", request.PeerInfo.UserID, request.Timestamp.UnixNano()))
	request.Signature = crypto.Sign(privateKey, payload)

	for _, peerAddr := range s.bootstrapPeers {
		if err := s.broadcaster(peerAddr, request); err != nil {
			return fmt.Errorf("failed to announce to peer %s: %w", peerAddr, err)
		}
	}
	return nil
}

// FindPeers retrieves a list of peers from the local peer store that are interested in a given frequency.
func (s *Service) FindPeers(frequency string) ([]types.PeerInfo, error) {
	allPeers := s.peerStore.List()
	matchingPeers := make([]types.PeerInfo, 0)

	for _, peer := range allPeers {
		for _, f := range peer.Frequencies {
			if f == frequency {
				matchingPeers = append(matchingPeers, *peer)
				break // Move to the next peer once a match is found
			}
		}
	}

	return matchingPeers, nil
}

// ProcessAnnouncement validates an incoming announcement and adds the peer to the store.
func (s *Service) ProcessAnnouncement(req *types.AnnounceRequest) error {
	if req == nil {
		return fmt.Errorf("invalid announcement request")
	}
	userID := req.PeerInfo.UserID
	if userID == "" {
		return fmt.Errorf("announcement from peer with no user ID")
	}

	user, err := s.identityService.GetUserByID(userID)
	if err != nil {
		return fmt.Errorf("failed to get user for announcement verification: %w", err)
	}

	payload := []byte(fmt.Sprintf("%s|%d", req.PeerInfo.UserID, req.Timestamp.UnixNano()))
	if !crypto.Verify(user.PublicKey, payload, req.Signature) {
		return fmt.Errorf("invalid signature on announcement from user %s", userID)
	}

	if err := s.peerStore.Add(&req.PeerInfo); err != nil {
		return fmt.Errorf("failed to add announced peer to store: %w", err)
	}
	return nil
}

// DiscoverPeers queries bootstrap peers to find other peers interested in a specific frequency.
func (s *Service) DiscoverPeers(frequency string) error {
	if s.requester == nil {
		return fmt.Errorf("requester is not configured")
	}

	request := &types.FindPeersRequest{
		Frequency: frequency,
	}

	for _, peerAddr := range s.bootstrapPeers {
		response, err := s.requester(peerAddr, request)
		if err != nil {
			return fmt.Errorf("failed to request peers from %s: %w", peerAddr, err)
		}

		for _, peer := range response.Peers {
			p := peer // Create a pointer to the peer to add to the store
			if err := s.peerStore.Add(&p); err != nil {
				fmt.Printf("Warning: failed to add discovered peer %s: %v\n", p.UserID, err)
			}
		}
	}
	return nil
}