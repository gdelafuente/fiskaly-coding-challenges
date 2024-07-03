package domain

import (
	"errors"
	"time"
)

var (
	ErrMissingSignatureID       = errors.New("missing signature id")
	ErrMissingSignatureDeviceID = errors.New("missing signature device id")
	ErrMissingSignatureRawData  = errors.New("missing signature raw data")
	ErrMissingSignatureValue    = errors.New("missing signature value")
	ErrMissingSignatureTime     = errors.New("missing signature time")
)

type Signature struct {
	deviceID  string
	id        string
	rawData   string
	value     []byte
	createdAt time.Time
}

func NewSignature(deviceID string, id string, rawData string, value []byte) (Signature, error) {
	s := Signature{
		deviceID:  deviceID,
		id:        id,
		rawData:   rawData,
		value:     value,
		createdAt: time.Now(),
	}

	return s, s.validate()
}

func (s Signature) validate() error {
	if s.id == "" {
		return ErrMissingSignatureID
	}
	if s.deviceID == "" {
		return ErrMissingSignatureDeviceID
	}
	if s.rawData == "" {
		return ErrMissingSignatureRawData
	}
	if len(s.value) == 0 {
		return ErrMissingSignatureValue
	}
	if s.createdAt.IsZero() {
		return ErrMissingSignatureTime
	}
	return nil
}

func (s Signature) DeviceID() string {
	return s.deviceID
}

func (s Signature) ID() string {
	return s.id
}

func (s Signature) RawData() string {
	return s.rawData
}

func (s Signature) Value() []byte {
	return s.value
}
