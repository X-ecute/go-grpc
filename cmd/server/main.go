package main

import (
	"log"

	"github.com/X-ecute/go-grpc/internal/db"
	"github.com/X-ecute/go-grpc/internal/rocket"
)

func Run() error {
	//responsible for initializing and starting the gRPC server
	rocketStore, err := db.New()
	if err != nil {
		return err
	}
	_ = rocket.New(rocketStore)
	return nil
}

func main() {
	if err := Run(); err != nil {
		log.Fatal(err)
	}
}
