package commands

import (
	"context"
	"errors"

	"github.com/gdelafuente/fiskaly-coding-challenges/signing-service-challenge-go/crypto"
	"github.com/gdelafuente/fiskaly-coding-challenges/signing-service-challenge-go/domain"
	"github.com/google/uuid"
)

var (
	ErrFetchingDevice    = errors.New("failed to fetch device")
	ErrBuildingSigner    = errors.New("failed to build signer")
	ErrSigning           = errors.New("failed to sign")
	ErrMissingDataToSign = errors.New("missing data to sign")
	ErrMissingDeviceID   = errors.New("missing device ID")
	ErrSignatureCreation = errors.New("failed to create a signature")
)

type createSignatureCommand struct {
	deviceID string
	data     string
}

func NewCreateSignatureCommand(deviceID string, data string) (createSignatureCommand, error) {
	cmd := createSignatureCommand{
		deviceID: deviceID,
		data:     data,
	}

	return cmd, cmd.validate()

}

func (c createSignatureCommand) validate() error {
	if c.data == "" {
		return errors.Join(ErrValidation, ErrMissingDataToSign)
	}
	if c.deviceID == "" {
		return errors.Join(ErrValidation, ErrMissingDeviceID)
	}
	return nil
}

type CreateSignatureCommandHandler struct {
	DeviceRepository      domain.DeviceRepository
	SignerFactoryResolver map[domain.SigningAlgorithm]crypto.SignerFactory
}

// TODO: this should return a DTO instead of a domain entity
func (h *CreateSignatureCommandHandler) Handle(ctx context.Context, cmd createSignatureCommand) (domain.Signature, error) {

	// Fail and retry in case of concurrent updates of the same device instead of locking
	var err error
	for retries := 0; retries < 3; retries++ {
		device, err := h.DeviceRepository.FindByID(ctx, cmd.deviceID)
		if err != nil {
			return domain.Signature{}, errors.Join(ErrFetchingDevice, err)
		}
		originalVersion := device.Version()
		enrichedData := device.EnrichData(cmd.data)

		signerFactory, ok := h.SignerFactoryResolver[device.Algorithm()]
		if !ok {
			return domain.Signature{}, ErrAlgorithmNotSupported
		}

		signer, err := signerFactory.Build(device.PrivateKey())
		if err != nil {
			return domain.Signature{}, errors.Join(ErrBuildingSigner, err)
		}

		signed, err := signer.Sign([]byte(enrichedData))
		if err != nil {
			return domain.Signature{}, errors.Join(ErrSigning, err)
		}

		signature, err := domain.NewSignature(device.ID(), uuid.NewString(), enrichedData, signed)
		if err != nil {
			return domain.Signature{}, errors.Join(ErrSignatureCreation, err)
		}
		device.AddSignature(signature)

		err = h.DeviceRepository.Update(ctx, device, originalVersion)
		if err == nil {
			return signature, nil
		}
		if err != domain.ErrDeviceNotFound {
			break
		}
	}

	return domain.Signature{}, err
}
