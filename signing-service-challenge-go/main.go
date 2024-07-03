package main

import (
	"log"
	"log/slog"

	"github.com/gdelafuente/fiskaly-coding-challenges/signing-service-challenge-go/api"
	"github.com/gdelafuente/fiskaly-coding-challenges/signing-service-challenge-go/application/commands"
	"github.com/gdelafuente/fiskaly-coding-challenges/signing-service-challenge-go/crypto"
	"github.com/gdelafuente/fiskaly-coding-challenges/signing-service-challenge-go/domain"
	"github.com/gdelafuente/fiskaly-coding-challenges/signing-service-challenge-go/persistence"
)

const (
	ListenAddress = ":8080"
)

func main() {

	logger := slog.Default()

	deviceRepository := persistence.NewInMemoryDeviceRepository()

	createDeviceCommandHandler := commands.CreateDeviceCommandHandler{
		DeviceRepository: deviceRepository,
		KeyProviderResolver: map[string]crypto.Provider{
			"rsa":   &crypto.RSAProvider{},
			"ecdsa": &crypto.ECDSAProvider{},
		},
	}

	createSignatureCommandHandler := commands.CreateSignatureCommandHandler{
		DeviceRepository: deviceRepository,
		SignerFactoryResolver: map[domain.SigningAlgorithm]crypto.SignerFactory{
			domain.SigningAlgorithmRSA:   &crypto.RSASignerFactory{},
			domain.SigningAlgorithmECDSA: &crypto.ECDSASignerFactory{},
		},
	}

	server := api.NewServer(ListenAddress, logger, createDeviceCommandHandler, createSignatureCommandHandler)

	if err := server.Run(); err != nil {
		log.Fatal("Could not start server on ", ListenAddress)
	}
}
