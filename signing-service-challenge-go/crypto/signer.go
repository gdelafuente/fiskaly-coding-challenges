package crypto

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}

type SignerFactory interface {
	Build(privateKey []byte) (Signer, error)
}

type RSASigner struct {
	privateKey []byte
	RSAMarshaler
}

func (s *RSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	keyPair, err := s.Unmarshal(s.privateKey)
	if err != nil {
		return nil, err
	}

	hashed := sha256.Sum256([]byte(dataToBeSigned))
	signature, err := rsa.SignPKCS1v15(
		rand.Reader,
		keyPair.Private,
		crypto.SHA256,
		hashed[:],
	)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

type RSASignerFactory struct {
}

func (f *RSASignerFactory) Build(privateKey []byte) (Signer, error) {
	return &RSASigner{
		privateKey:   privateKey,
		RSAMarshaler: RSAMarshaler{},
	}, nil
}

type ECDSASigner struct {
	privateKey []byte
	ECCMarshaler
}

func (s *ECDSASigner) Sign(dataToBeSigned []byte) ([]byte, error) {
	keyPair, err := s.Decode(s.privateKey)
	if err != nil {
		return nil, err
	}

	signature, err := ecdsa.SignASN1(
		rand.Reader,
		keyPair.Private,
		dataToBeSigned,
	)
	if err != nil {
		return nil, err
	}

	return signature, nil
}

type ECDSASignerFactory struct {
}

func (f *ECDSASignerFactory) Build(privateKey []byte) (Signer, error) {
	return &ECDSASigner{
		privateKey:   privateKey,
		ECCMarshaler: ECCMarshaler{},
	}, nil
}
