package api

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

type NewDeviceReq struct {
	ID        string `json:"device_id"`
	Algorithm string `json:"algorithm"`
	Label     string `json:"label"`
}

// CreateSignatureDevice(id: string, algorithm: 'ECC' | 'RSA', [optional]: label: string): CreateSignatureDeviceResponse`
func (s *Server) CreateSignatureDevice(response http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var dev NewDeviceReq
	err := decoder.Decode(&dev)
	if err != nil {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			err.Error(),
		})
	}
	ctx := request.Context()
	err = s.singin.CreateSigner(ctx, dev.ID, dev.Algorithm, dev.Label)
	if err != nil {
		slog.Error("CreateSigner", "error", fmt.Sprintf("%+v", err))
		WriteInternalError(response)
	}
	WriteAPIResponse(response, http.StatusCreated, nil)
}

func (s *Server) ListDevices(response http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	list, err := s.singin.ListDevices(ctx)
	if err != nil {
		slog.Error("ListDevices", "error", fmt.Sprintf("%+v", err))
		WriteInternalError(response)
	}
	WriteAPIResponse(response, http.StatusOK, DevListToHTTP(list))
}

func (s *Server) DeleteDevice(response http.ResponseWriter, request *http.Request) {
	idString := request.PathValue("id")
	if idString == "" {
		WriteErrorResponse(response, http.StatusBadRequest, nil)
	}
	ctx := request.Context()
	err := s.singin.DeleteSigner(ctx, idString)
	if err != nil {
		slog.Error("DeleteSigner", "error", fmt.Sprintf("%+v", err))
		WriteInternalError(response)
	}
}

type APIDevice struct {
	ID               string `json:"id"`
	Algorithm        string `json:"algrithm"`
	SignatureCounter int    `json:"signature_counter"`
	Label            string `json:"label"`
	LastSignature    string `json:"last_signature"`
}

func DevListToHTTP(list []domain.Device) []APIDevice {
	ret := make([]APIDevice, 0, len(list))
	for _, dev := range list {
		ret = append(ret, APIDevice{
			ID:               dev.ID,
			Algorithm:        dev.Algorithm.String(),
			SignatureCounter: dev.SignatureCounter,
			Label:            dev.Label,
			LastSignature:    dev.LastSignature,
		})
	}
	return ret
}
