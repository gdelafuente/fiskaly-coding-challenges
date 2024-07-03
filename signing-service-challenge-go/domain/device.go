package domain

import (
	"context"
	"encoding/base64"
	"errors"
	"strconv"
)

var (
	ErrMissingDeviceID         = errors.New("missing device id")
	ErrMissingDevicePublicKey  = errors.New("missing device public key")
	ErrMissingDevicePrivateKey = errors.New("missing device private key")
)

type Device struct {
	id               string
	signingAlgorithm SigningAlgorithm
	publicKey        []byte
	privateKey       []byte
	label            string
	version          int
	signatures       []Signature
}

// TODO: having a list with all the signatures means we'll be loading all of them
// even if we don't need them, e.g. when creating a new signature.
// I'd turn Signature into an aggregate root with it's own repositories or views,
// which would support filtering, ordering and pagination
func NewDevice(id, algorithmName, label string, publicKey, privateKey []byte) (Device, error) {
	Algorithm, err := NewSigningAlgorithm(algorithmName)
	if err != nil {
		return Device{}, err
	}

	d := Device{
		id:               id,
		signingAlgorithm: Algorithm,
		publicKey:        publicKey,
		privateKey:       privateKey,
		label:            label,
		version:          0,
	}

	return d, d.validate()
}

func (d Device) validate() error {
	if d.id == "" {
		return ErrMissingDeviceID
	}
	if len(d.publicKey) == 0 {
		return ErrMissingDevicePublicKey
	}
	if len(d.privateKey) == 0 {
		return ErrMissingDevicePrivateKey
	}
	return nil
}

func (d Device) EnrichData(data string) string {
	data = strconv.Itoa(d.SignaturesCount()) + "_" + data + "_"
	if d.SignaturesCount() == 0 {
		data += base64.StdEncoding.EncodeToString([]byte(d.ID()))
	} else {
		latestSignature := d.signatures[len(d.signatures)-1]
		data += base64.StdEncoding.EncodeToString(latestSignature.Value())
	}
	return data
}

func (d Device) ID() string {
	return d.id
}

func (d Device) Algorithm() SigningAlgorithm {
	return d.signingAlgorithm
}

func (d Device) PrivateKey() []byte {
	return d.privateKey
}

func (d Device) PublicKey() []byte {
	return d.publicKey
}

func (d Device) Label() string {
	return d.label
}

func (d Device) Version() int {
	return d.version
}

func (d Device) SignaturesCount() int {
	return len(d.signatures)
}

func (d Device) Signatures() []Signature {
	return d.signatures
}

func (d *Device) AddSignature(signature Signature) {
	d.signatures = append(d.signatures, signature)
	d.version++
}

var (
	ErrDeviceNotFound        = errors.New("device not found")
	ErrDeviceVersionMismatch = errors.New("device version mismatch")
)

type DeviceRepository interface {
	Save(ctx context.Context, d Device) error
	Update(ctx context.Context, d Device, expectedVersion int) error
	FindByID(ctx context.Context, id string) (Device, error)
	ListAll(ctx context.Context) ([]Device, error)
}
