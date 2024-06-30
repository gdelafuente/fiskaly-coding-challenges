package main

import (
	"log"

	"github.com/gdelafuente/fiskaly-coding-challenges/signing-service-challenge-go/api"
)

const (
	ListenAddress = ":8080"
	// TODO: add further configuration parameters here ...
)

func main() {
	server := api.NewServer(ListenAddress)

	if err := server.Run(); err != nil {
		log.Fatal("Could not start server on ", ListenAddress)
	}
}
