package commands

import (
	"context"
	"errors"

	"github.com/gdelafuente/fiskaly-coding-challenges/signing-service-challenge-go/crypto"
	"github.com/gdelafuente/fiskaly-coding-challenges/signing-service-challenge-go/domain"
	"github.com/google/uuid"
)

var (
	ErrDeviceIDAlreadyInUse = errors.New("device id already in use")
	ErrKeyGeneration        = errors.New("failed to generate keys")
	ErrDeviceCreation       = errors.New("failed to create a device")
	ErrMissingAlgorithmName = errors.New("missing algorithm name")
)

type createDeviceCommand struct {
	algorithmName string
	label         string
}

func NewCreateDeviceCommand(algorithmName string, label string) (createDeviceCommand, error) {
	cmd := createDeviceCommand{
		algorithmName: algorithmName,
		label:         label,
	}
	return cmd, cmd.validate()
}

func (c createDeviceCommand) validate() error {
	if c.algorithmName == "" {
		return errors.Join(ErrValidation, ErrMissingAlgorithmName)
	}
	return nil
}

type CreateDeviceCommandHandler struct {
	DeviceRepository    domain.DeviceRepository
	KeyProviderResolver map[string]crypto.Provider
}

// TODO: this should return a DTO instead of a domain entity
func (h *CreateDeviceCommandHandler) Handle(ctx context.Context, cmd createDeviceCommand) (domain.Device, error) {
	id := uuid.NewString()

	_, err := h.DeviceRepository.FindByID(ctx, id)
	if err == nil {
		return domain.Device{}, ErrDeviceIDAlreadyInUse
	}

	keyProvider, ok := h.KeyProviderResolver[cmd.algorithmName]
	if !ok {
		return domain.Device{}, errors.Join(ErrValidation, ErrAlgorithmNotSupported)
	}

	keyPair, err := keyProvider.Provide()
	if err != nil {
		return domain.Device{}, errors.Join(ErrKeyGeneration, err)
	}

	device, err := domain.NewDevice(id, cmd.algorithmName, cmd.label, keyPair.Public, keyPair.Private)
	if err != nil {
		return domain.Device{}, errors.Join(ErrDeviceCreation, err)
	}

	err = h.DeviceRepository.Save(ctx, device)
	if err != nil {
		return domain.Device{}, errors.Join(ErrSavingDevice, err)
	}

	return device, nil
}
