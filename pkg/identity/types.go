package identity

import (
	"time"
)

// NexusUserObjectV1 represents a sovereign user identity on the network.
// This is a native Go struct.
type NexusUserObjectV1 struct {
	UserId       string
	PublicKey    []byte
	Handle       string
	CreatedAt    time.Time
	MetadataURI  string
}