package api

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gdelafuente/fiskaly-coding-challenges/signing-service-challenge-go/application/commands"
)

func (s *Server) Devices(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.CreateDevice(w, r)
	default:
		WriteErrorResponse(w, http.StatusMethodNotAllowed, []string{
			http.StatusText(http.StatusMethodNotAllowed),
		})
	}
}

type CreateDeviceRequest struct {
	Algorithm string `json:"algorithm"`
	Label     string `json:"label"`
}

type DeviceResponse struct {
	ID              string `json:"id"`
	Algorithm       string `json:"algorithm"`
	Label           string `json:"label"`
	PublicKey       []byte `json:"public_key"`
	SignaturesCount int    `json:"signatures_count"`
}

func (s *Server) CreateDevice(w http.ResponseWriter, r *http.Request) {
	var request CreateDeviceRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		s.logger.Info("Invalid device creation request", slog.String("error", err.Error()))
		WriteErrorResponse(w, http.StatusBadRequest, []string{
			http.StatusText(http.StatusBadRequest),
			err.Error(),
		})
		return
	}

	cmd, err := commands.NewCreateDeviceCommand(request.Algorithm, request.Label)
	if err != nil {
		s.logger.Info("Invalid device creation command", slog.String("error", err.Error()))
		WriteErrorResponse(w, http.StatusBadRequest, []string{
			http.StatusText(http.StatusBadRequest),
			err.Error(),
		})
		return
	}

	device, err := s.createDeviceCommandHandler.Handle(context.Background(), cmd)
	if err != nil {
		if errors.Is(err, commands.ErrValidation) {
			s.logger.Info("Invalid device creation command", slog.String("error", err.Error()))
			WriteErrorResponse(w, http.StatusBadRequest, []string{
				http.StatusText(http.StatusBadRequest),
				err.Error(),
			})
			return
		}
		s.logger.Error("Failed to create a device", slog.String("error", err.Error()))
		WriteErrorResponse(w, http.StatusInternalServerError, []string{
			// Avoid propagating internal errors traces to the clients
			http.StatusText(http.StatusInternalServerError),
		})
		return
	}

	response := DeviceResponse{
		ID:              device.ID(),
		Algorithm:       string(device.Algorithm()),
		Label:           device.Label(),
		PublicKey:       device.PublicKey(),
		SignaturesCount: device.SignaturesCount(),
	}
	WriteAPIResponse(w, http.StatusOK, response)
}
