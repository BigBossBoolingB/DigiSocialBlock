package discovery

import (
	"crypto/ed25519"
	"fmt"
	"time"

	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/crypto"
	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/echonet/peerstore"
	pb "github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/proto/echonet/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Broadcaster defines the interface for sending an announcement to a specific peer.
type Broadcaster func(peerAddress string, request *pb.AnnounceRequest) error

// Requester defines the interface for sending a request to a peer and getting a response.
type Requester func(peerAddress string, request *pb.FindPeersRequest) (*pb.FindPeersResponse, error)

// Service handles the peer discovery process.
type Service struct {
	peerStore      *peerstore.PeerStore
	bootstrapPeers []string
	broadcaster    Broadcaster
	requester      Requester
}

// NewService creates a new discovery service.
func NewService(ps *peerstore.PeerStore, bootstrapPeers []string, broadcaster Broadcaster, requester Requester) *Service {
	return &Service{
		peerStore:      ps,
		bootstrapPeers: bootstrapPeers,
		broadcaster:    broadcaster,
		requester:      requester,
	}
}

// Announce broadcasts the presence of the local node to the network.
// It signs the node's PeerInfo and sends it to all known bootstrap peers.
func (s *Service) Announce(localPeer *pb.PeerInfo, privateKey ed25519.PrivateKey) error {
	if localPeer == nil || localPeer.UserId == "" {
		return fmt.Errorf("cannot announce with invalid peer info")
	}

	// 1. Create the core announcement request
	request := &pb.AnnounceRequest{
		PeerInfo:  localPeer,
		Timestamp: timestamppb.New(time.Now()),
	}

	// 2. Serialize the request payload for signing
	// We sign the PeerInfo and Timestamp together.
	payload, err := proto.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal announcement request for signing: %w", err)
	}

	// 3. Sign the payload
	signature := crypto.Sign(privateKey, payload)
	request.Signature = signature

	// 4. Broadcast the signed request to all bootstrap peers
	for _, peerAddr := range s.bootstrapPeers {
		if err := s.broadcaster(peerAddr, request); err != nil {
			// In a real implementation, we might collect errors and continue,
			// but for now, we'll return on the first error.
			return fmt.Errorf("failed to announce to peer %s: %w", peerAddr, err)
		}
	}

	return nil
}

// DiscoverPeers queries bootstrap peers to find other peers interested in a specific frequency.
func (s *Service) DiscoverPeers(frequency string) error {
	if s.requester == nil {
		return fmt.Errorf("requester is not configured")
	}

	request := &pb.FindPeersRequest{
		Frequency: frequency,
	}

	for _, peerAddr := range s.bootstrapPeers {
		response, err := s.requester(peerAddr, request)
		if err != nil {
			// In a real implementation, we might log the error and try the next peer.
			return fmt.Errorf("failed to request peers from %s: %w", peerAddr, err)
		}

		for _, peer := range response.Peers {
			if err := s.peerStore.Add(peer); err != nil {
				// Log error but continue processing other peers
				fmt.Printf("Warning: failed to add discovered peer %s: %v\n", peer.UserId, err)
			}
		}
	}

	return nil
}