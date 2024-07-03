package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gdelafuente/fiskaly-coding-challenges/signing-service-challenge-go/application/commands"
	"github.com/go-chi/chi"
)

// Response is the generic API response container.
type Response struct {
	Data interface{} `json:"data"`
}

// ErrorResponse is the generic error API response container.
type ErrorResponse struct {
	Errors []string `json:"errors"`
}

// Server manages HTTP requests and dispatches them to the appropriate services.
type Server struct {
	listenAddress                 string
	logger                        *slog.Logger
	createDeviceCommandHandler    commands.CreateDeviceCommandHandler
	createSignatureCommandHandler commands.CreateSignatureCommandHandler
}

// NewServer is a factory to instantiate a new Server.
func NewServer(listenAddress string, logger *slog.Logger, createDeviceCommandHandler commands.CreateDeviceCommandHandler, createSignatureCommandHandler commands.CreateSignatureCommandHandler) *Server {
	return &Server{
		listenAddress:                 listenAddress,
		logger:                        logger,
		createDeviceCommandHandler:    createDeviceCommandHandler,
		createSignatureCommandHandler: createSignatureCommandHandler,
	}
}

// Run registers all HandlerFuncs for the existing HTTP routes and starts the Server.
func (s *Server) Run() error {
	router := chi.NewRouter()
	router.Route("/api/v0", func(r chi.Router) {
		r.Handle("/health", http.HandlerFunc(s.Health))
		r.Handle("/devices", http.HandlerFunc(s.Devices))
		r.Handle("/devices/{deviceID}/signatures", http.HandlerFunc(s.Signatures))
	})

	s.logger.Info(fmt.Sprintf("Starting HTTP server listening on %s", s.listenAddress))
	return http.ListenAndServe(s.listenAddress, router)
}

// WriteInternalError writes a default internal error message as an HTTP response.
func WriteInternalError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
}

// WriteErrorResponse takes an HTTP status code and a slice of errors
// and writes those as an HTTP error response in a structured format.
func WriteErrorResponse(w http.ResponseWriter, code int, errors []string) {
	w.WriteHeader(code)

	errorResponse := ErrorResponse{
		Errors: errors,
	}

	bytes, err := json.Marshal(errorResponse)
	if err != nil {
		WriteInternalError(w)
	}

	w.Write(bytes)
}

// WriteAPIResponse takes an HTTP status code and a generic data struct
// and writes those as an HTTP response in a structured format.
func WriteAPIResponse(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)

	response := Response{
		Data: data,
	}

	bytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		WriteInternalError(w)
	}

	w.Write(bytes)
}
