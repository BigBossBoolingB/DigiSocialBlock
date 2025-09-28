package echonet

import (
	"log"
	"net"
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

// Server wraps the Go RPC server and our API implementation.
type Server struct {
	api      *EchoNetAPI
	listener net.Listener
}

// NewServer creates a new EchoNet RPC server.
func NewServer(ds *discovery.Service) (*Server, error) {
	api := NewEchoNetAPI(ds)
	err := rpc.Register(api)
	if err != nil {
		// Ignore "service already defined" error which can happen during hot-reloading in tests.
		if err.Error() != "rpc: service already defined: EchoNetAPI" {
			return nil, err
		}
	}
	rpc.HandleHTTP() // Use default HTTP handlers

	return &Server{
		api: api,
	}, nil
}

// Start begins listening for RPC requests on the given address.
func (s *Server) Start(address string) error {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	s.listener = l
	log.Printf("EchoNet RPC server listening on %s", s.Address())
	go rpc.Accept(l)
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