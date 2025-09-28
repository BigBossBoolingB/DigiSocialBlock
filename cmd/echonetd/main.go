package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/content"
	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/echonet"
	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/echonet/discovery"
	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/echonet/peerstore"
	"github.com/BigBossBoolingB/Digital-Golem-Engine/pkg/identity"
)

func main() {
	// 1. Setup command-line flags
	listenAddr := flag.String("listen", ":8080", "Address for the RPC server to listen on")
	flag.Parse()

	// 2. Initialize all core services
	fmt.Println("Initializing services...")
	idService := identity.NewService()
	contentService := content.NewService(idService)
	ps := peerstore.New()

	// A real application would get these from a config file or a different discovery mechanism.
	bootstrapPeers := []string{} // No bootstrap peers for a single node running in isolation.

	// Wire in the live networking functions
	ds := discovery.NewService(idService, ps, bootstrapPeers, discovery.LiveAnnounceBroadcaster, discovery.LiveFindPeersRequester)

	// 3. Create and start the EchoNet RPC Server
	server, err := echonet.NewServer(ds, contentService)
	if err != nil {
		fmt.Printf("Error creating server: %v\n", err)
		os.Exit(1)
	}

	go func() {
		if err := server.Start(*listenAddr); err != nil {
			fmt.Printf("Error starting server: %v\n", err)
			os.Exit(1)
		}
	}()

	fmt.Printf("EchoNet Daemon started. Listening on %s\n", *listenAddr)
	fmt.Println("Press Ctrl+C to exit.")

	// 4. Wait for a shutdown signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	// 5. Graceful shutdown
	fmt.Println("\nShutting down EchoNet Daemon...")
	server.Stop()
	fmt.Println("Server stopped.")
}