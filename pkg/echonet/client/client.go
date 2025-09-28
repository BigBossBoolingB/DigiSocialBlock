package client

import (
	"net/rpc"

	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/echonet/types"
)

// Client is a wrapper for a Go net/rpc client for our EchoNet service.
type Client struct {
	rpcClient *rpc.Client
}

// New creates a new RPC client connected to the given server address.
func New(serverAddress string) (*Client, error) {
	rpcClient, err := rpc.DialHTTP("tcp", serverAddress)
	if err != nil {
		return nil, err
	}
	return &Client{
		rpcClient: rpcClient,
	}, nil
}

// Announce sends an announcement to the connected peer.
func (c *Client) Announce(args *types.AnnounceRequest) (*types.AnnounceResponse, error) {
	var reply types.AnnounceResponse
	err := c.rpcClient.Call("EchoNetAPI.Announce", args, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

// FindPeers requests a list of peers for a given frequency from the connected peer.
func (c *Client) FindPeers(args *types.FindPeersRequest) (*types.FindPeersResponse, error) {
	var reply types.FindPeersResponse
	err := c.rpcClient.Call("EchoNetAPI.FindPeers", args, &reply)
	if err != nil {
		return nil, err
	}
	return &reply, nil
}

// Close closes the underlying RPC connection.
func (c *Client) Close() error {
	return c.rpcClient.Close()
}