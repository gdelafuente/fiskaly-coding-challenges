package domain_test

import (
	"encoding/base64"
	"errors"
	"testing"

	"github.com/gdelafuente/fiskaly-coding-challenges/signing-service-challenge-go/domain"
)

func Test_NewDevice_OK(t *testing.T) {
	id := "device_id_0"
	privateKey := []byte("private_key_0")
	publicKey := []byte("public_key_0")
	label := "device_label_0"
	device, err := domain.NewDevice(id, "rsa", label, publicKey, privateKey)

	if err != nil {
		t.Fatal("Expected no error, got", err)
	}

	if device.ID() != id {
		t.Fatal("Expected id to be", id, "got", device.ID())
	}
	if device.Algorithm() != domain.SigningAlgorithmRSA {
		t.Fatal("Expected algorithm to be rsa, got", device.Algorithm())
	}
	if string(device.PrivateKey()) != string(privateKey) {
		t.Fatal("Expected private key to be", string(privateKey), "got", string(device.PrivateKey()))
	}
	if string(device.PublicKey()) != string(publicKey) {
		t.Fatal("Expected public key to be", string(publicKey), "got", string(device.PublicKey()))
	}
	if device.Label() != label {
		t.Fatal("Expected label to be", label, "got", device.Label())
	}
}

func Test_NewDevice_EmptyID_Error(t *testing.T) {
	_, err := domain.NewDevice("", "rsa", "device_label_0", []byte("public_key_0"), []byte("private_key_0"))

	expectedError := domain.ErrMissingDeviceID
	if err == nil || !errors.Is(err, expectedError) {
		t.Fatal("Expected error to be", expectedError, "got", err)
	}
}

func Test_NewDevice_EmptyPublicKey_Error(t *testing.T) {
	_, err := domain.NewDevice("device_id_0", "rsa", "device_label_0", []byte{}, []byte("private_key_0"))

	expectedError := domain.ErrMissingDevicePublicKey
	if err == nil || !errors.Is(err, expectedError) {
		t.Fatal("Expected error to be", expectedError, "got", err)
	}
}

func Test_NewDevice_EmptyPrivateKey_Error(t *testing.T) {
	_, err := domain.NewDevice("device_id_0", "rsa", "device_label_0", []byte("public_key"), []byte{})

	expectedError := domain.ErrMissingDevicePrivateKey
	if err == nil || !errors.Is(err, expectedError) {
		t.Fatal("Expected error to be", expectedError, "got", err)
	}
}

func Test_NewDevice_InvalidAlgorithm_Error(t *testing.T) {
	_, err := domain.NewDevice("device_id_0", "invalid_algorithm", "device_label_0", []byte("public_key_0"), []byte("private_key_0"))

	expectedError := domain.ErrUnknownSigningAlgorithm
	if err == nil || !errors.Is(err, expectedError) {
		t.Fatal("Expected error to be", expectedError, "got", err)
	}
}

func Test_Device_AddSignature(t *testing.T) {
	device, err := domain.NewDevice("device_id_0", "rsa", "device_label_0", []byte("public_key_0"), []byte("private_key_0"))
	if err != nil {
		t.Fatal("Expected no error, got", err)
	}

	signature, err := domain.NewSignature("device_id_0", "signature_id_0", "foo", []byte("signature_0"))
	if err != nil {
		t.Fatal("Expected no error, got", err)
	}

	previousSignatureCount := device.SignaturesCount()
	device.AddSignature(signature)

	expectedSignatureCount := previousSignatureCount + 1
	if device.SignaturesCount() != expectedSignatureCount {
		t.Fatal("Expected signature count to be", expectedSignatureCount, "got", device.SignaturesCount())
	}

	if device.Signatures()[0].ID() != signature.ID() {
		t.Fatal("Expected signature to be", signature, "got", device.Signatures()[0])
	}

}

func Test_Device_EnrichData(t *testing.T) {
	device, err := domain.NewDevice("device_id_0", "rsa", "device_label_0", []byte("public_key_0"), []byte("private_key_0"))
	if err != nil {
		t.Fatal("Expected no error, got", err)
	}

	signature, err := domain.NewSignature("device_id_0", "signature_id_0", "foo", []byte("signature_0"))
	if err != nil {
		t.Fatal("Expected no error, got", err)
	}
	device.AddSignature(signature)

	dataToBeSigned := "data_to_be_signed"

	enrichedData := device.EnrichData(dataToBeSigned)

	expectedEnrichedData := "1" + "_" + dataToBeSigned + "_" + base64.StdEncoding.EncodeToString(signature.Value())
	if enrichedData != expectedEnrichedData {
		t.Fatal("Expected enriched data to be", expectedEnrichedData, "got", enrichedData)
	}

}

func Test_Device_EnrichData_FirstSignature(t *testing.T) {
	device, err := domain.NewDevice("device_id_0", "rsa", "device_label_0", []byte("public_key_0"), []byte("private_key_0"))
	if err != nil {
		t.Fatal("Expected no error, got", err)
	}

	dataToBeSigned := "data_to_be_signed"

	enrichedData := device.EnrichData(dataToBeSigned)

	expectedEnrichedData := "0" + "_" + dataToBeSigned + "_" + base64.StdEncoding.EncodeToString([]byte(device.ID()))
	if enrichedData != expectedEnrichedData {
		t.Fatal("Expected enriched data to be", expectedEnrichedData, "got", enrichedData)
	}

}
