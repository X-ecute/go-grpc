package main

import (
	"log"

	"github.com/X-ecute/go-grpc/internal/db"
	"github.com/X-ecute/go-grpc/internal/rocket"
	"github.com/X-ecute/go-grpc/internal/transport/grpc"
)

func Run() error {
	// responsible for initializing and starting the gRPC server
	rocketStore, err := db.New()
	if err != nil {
		return err
	}
	err = rocketStore.Migrate()
	if err != nil {
		log.Fatal("failed to migrate rocket store")
		return err
	}

	// Create rocket service
	rktService := rocket.New(&rocketStore) // Remove & since New() should return *Store

	// Create gRPC handler
	rktHandler := grpc.New(&rktService) // Remove & since New() returns *Handler

	if err := rktHandler.Serve(); err != nil {
		return err
	}
	return nil
}

func main() {
	if err := Run(); err != nil {
		log.Fatal(err)
	}
}
