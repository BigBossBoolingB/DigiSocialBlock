package discovery

import (
	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/echonet/client"
	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/echonet/types"
)

// LiveAnnounceBroadcaster is a concrete implementation of the Broadcaster func type.
// It creates a new RPC client for each call, sends the announcement, and then closes the connection.
func LiveAnnounceBroadcaster(peerAddress string, request *types.AnnounceRequest) error {
	c, err := client.New(peerAddress)
	if err != nil {
		return err
	}
	defer c.Close()

	_, err = c.Announce(request)
	return err
}

// LiveFindPeersRequester is a concrete implementation of the Requester func type.
// It creates a new RPC client for each call, requests peers, and then closes the connection.
func LiveFindPeersRequester(peerAddress string, request *types.FindPeersRequest) (*types.FindPeersResponse, error) {
	c, err := client.New(peerAddress)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	return c.FindPeers(request)
}