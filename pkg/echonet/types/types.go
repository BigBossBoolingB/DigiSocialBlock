package types

import (
	"time"
)

// PeerInfo contains the necessary information for a node to connect to a peer.
type PeerInfo struct {
	UserID       string
	Multiaddress string
	Frequencies  []string
}

// AnnounceRequest defines the arguments for the Announce RPC method.
type AnnounceRequest struct {
	PeerInfo  PeerInfo
	Timestamp time.Time
	Signature []byte
}

// AnnounceResponse defines the response for the Announce RPC method.
type AnnounceResponse struct {
	ConfirmationMessage string
}

// FindPeersRequest defines the arguments for the FindPeers RPC method.
type FindPeersRequest struct {
	Frequency string
}

// FindPeersResponse defines the response for the FindPeers RPC method.
type FindPeersResponse struct {
	Peers []PeerInfo
}