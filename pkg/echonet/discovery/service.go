package discovery

import (
	"crypto/ed25519"
	"fmt"
	"time"

	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/crypto"
	pb "github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/proto/echonet/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Broadcaster defines the interface for sending a message to a specific peer.
// This allows us to mock the network layer for testing.
type Broadcaster func(peerAddress string, request *pb.AnnounceRequest) error

// Service handles the peer discovery process.
type Service struct {
	bootstrapPeers []string
	broadcaster    Broadcaster
}

// NewService creates a new discovery service.
func NewService(bootstrapPeers []string, broadcaster Broadcaster) *Service {
	return &Service{
		bootstrapPeers: bootstrapPeers,
		broadcaster:    broadcaster,
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