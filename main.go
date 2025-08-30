// cmd/api/main.go
package main

import (
	"log"

	"github.com/lm-Alesh-Patil/notification-api-service/config"
	"github.com/lm-Alesh-Patil/notification-api-service/server"
)

func main() {
	// Load config
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Create server
	s := &server.Server{Config: cfg}

	// Setup DB, Redis, etc.
	if err := s.Setup(); err != nil {
		log.Fatalf("server setup failed: %v", err)
	}

	// Start HTTP server
	if err := s.Start(); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
