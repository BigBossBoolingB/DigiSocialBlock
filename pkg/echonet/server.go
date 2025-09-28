package echonet

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/echonet/discovery"
	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/echonet/types"
)

// EchoNetAPI is the struct that will be registered as our RPC service.
type EchoNetAPI struct {
	Discovery *discovery.Service
}

// NewEchoNetAPI creates a new API handler.
func NewEchoNetAPI(ds *discovery.Service) *EchoNetAPI {
	return &EchoNetAPI{
		Discovery: ds,
	}
}

// Announce is the RPC method that allows a peer to announce its presence.
func (e *EchoNetAPI) Announce(args *types.AnnounceRequest, reply *types.AnnounceResponse) error {
	if err := e.Discovery.ProcessAnnouncement(args); err != nil {
		return err
	}
	*reply = types.AnnounceResponse{
		ConfirmationMessage: "Announcement for " + args.PeerInfo.UserID + " processed.",
	}
	return nil
}

// FindPeers is the RPC method that allows a peer to find other peers by frequency.
func (e *EchoNetAPI) FindPeers(args *types.FindPeersRequest, reply *types.FindPeersResponse) error {
	peers, err := e.Discovery.FindPeers(args.Frequency)
	if err != nil {
		return err
	}
	reply.Peers = peers
	return nil
}

// Server wraps the Go RPC server and our API implementation.
type Server struct {
	rpcServer *rpc.Server
	listener  net.Listener
}

// NewServer creates a new EchoNet RPC server that is isolated and safe for parallel tests.
func NewServer(ds *discovery.Service) (*Server, error) {
	api := NewEchoNetAPI(ds)
	rpcServer := rpc.NewServer() // Create a new, isolated RPC server instance
	err := rpcServer.Register(api)
	if err != nil {
		return nil, err
	}

	return &Server{
		rpcServer: rpcServer,
	}, nil
}

// Start begins listening for HTTP RPC requests on the given address.
func (s *Server) Start(address string) error {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	s.listener = l
	log.Printf("EchoNet RPC server listening on %s", s.Address())

	// Serve the isolated RPC server over HTTP. This is safe for parallel tests.
	go http.Serve(l, s.rpcServer)
	return nil
}

// Stop closes the server's listener.
func (s *Server) Stop() {
	if s.listener != nil {
		s.listener.Close()
	}
}

// Address returns the address the server is listening on.
func (s *Server) Address() string {
	if s.listener != nil {
		return s.listener.Addr().String()
	}
	return ""
}