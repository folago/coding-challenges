package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/service"
)

// Response is the generic API response container.
type Response struct {
	Data interface{} `json:"data,omitempty"`
}

// ErrorResponse is the generic error API response container.
type ErrorResponse struct {
	Errors []string `json:"errors"`
}

// Server manages HTTP requests and dispatches them to the appropriate services.
type Server struct {
	listenAddress string
	singin        service.Signer
}

// NewServer is a factory to instantiate a new Server.
func NewServer(listenAddress string) *Server {

	return &Server{
		listenAddress: listenAddress,
		singin: service.NewSignerService(
			persistence.NewMemDeviceRepository(),
			persistence.NewMemSignatureRepository()),
	}
}

// Run registers all HandlerFuncs for the existing HTTP routes and starts the Server.
// TODO: add more middlewares, at least a recovery one
func (s *Server) Run() error {
	mux := http.NewServeMux()

	mux.Handle("GET /api/v0/health", logs(s.Health))

	mux.Handle("POST /api/v0/devices", logs(s.CreateSignatureDevice))
	mux.Handle("GET /api/v0/devices", logs(s.ListDevices))
	mux.Handle("DELETE /api/v0/devices/{id}", logs(s.DeleteDevice))

	mux.Handle("POST /api/v0/signatures", logs(s.SignTransaction))
	mux.Handle("GET /api/v0/signatures", logs(s.ListTransactions))

	slog.Info("server running", "address", s.listenAddress)
	return http.ListenAndServe(s.listenAddress, mux)
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

func logs(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slog.Info("request", "method", r.Method, "url", r.URL, "ts", time.Now())
		next(w, r)
	}
}
