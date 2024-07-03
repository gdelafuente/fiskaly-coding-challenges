package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gdelafuente/fiskaly-coding-challenges/signing-service-challenge-go/application/commands"
	"github.com/gdelafuente/fiskaly-coding-challenges/signing-service-challenge-go/domain"
)

func (s *Server) Signatures(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.CreateDeviceSignature(w, r)
	default:
		WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
	}
}

type CreateDeviceSignatureRequest struct {
	Data string `json:"data"`
}

type SignatureResponse struct {
	DeviceID   string `json:"device_id"`
	ID         string `json:"id"`
	Signature  []byte `json:"signature"`
	SignedData string `json:"signed_data"`
}

func (s *Server) CreateDeviceSignature(w http.ResponseWriter, r *http.Request) {
	var request CreateDeviceSignatureRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		s.logger.Info("Invalid signature creation request", slog.String("error", err.Error()))
		WriteErrorResponse(w, http.StatusBadRequest, []string{
			http.StatusText(http.StatusBadRequest),
			err.Error(),
		})
		return
	}

	deviceID := strings.Split(strings.Split(r.URL.Path, "/devices/")[1], "/signatures")[0]
	cmd, err := commands.NewCreateSignatureCommand(deviceID, request.Data)
	if err != nil {
		s.logger.Info("Invalid signature creation command", slog.String("error", err.Error()))
		WriteErrorResponse(w, http.StatusBadRequest, []string{
			http.StatusText(http.StatusBadRequest),
			err.Error(),
		})
		return
	}

	signature, err := s.createSignatureCommandHandler.Handle(context.Background(), cmd)
	if err != nil {
		if errors.Is(err, domain.ErrDeviceNotFound) {
			s.logger.Info("Device for creating a signature not found", slog.String("error", err.Error()))
			WriteErrorResponse(w, http.StatusNotFound, []string{
				http.StatusText(http.StatusNotFound),
			})
			return
		}
		s.logger.Error("Failed to create a signature", slog.String("error", err.Error()))
		WriteErrorResponse(w, http.StatusInternalServerError, []string{
			// Avoid propagating internal errors traces to the clients
			http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	response := SignatureResponse{
		DeviceID:   signature.DeviceID(),
		ID:         signature.ID(),
		Signature:  signature.Value(),
		SignedData: signature.RawData(),
	}
	WriteAPIResponse(w, http.StatusOK, response)
}
