package domain

import "errors"

type SigningAlgorithm string

var (
	ErrUnknownSigningAlgorithm = errors.New("unknown signing algorithm")
)

const (
	SigningAlgorithmRSA   SigningAlgorithm = "rsa"
	SigningAlgorithmECDSA SigningAlgorithm = "ecdsa"
)

func (s SigningAlgorithm) validate() error {
	switch s {
	case SigningAlgorithmRSA:
		return nil
	case SigningAlgorithmECDSA:
		return nil
	}
	return ErrUnknownSigningAlgorithm
}

func NewSigningAlgorithm(val string) (SigningAlgorithm, error) {
	s := SigningAlgorithm(val)
	return s, s.validate()
}
