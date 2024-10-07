package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

// SignTransactionReq represent the request for signing some transaction data.
// Let's assume the data is base64 encoded to fit in a JSON document.
type SignTransactionReq struct {
	ID   string `json:"device_id"`
	Data string `json:"data"`
}

// SignTransaction(deviceId: string, data: string): SignatureResponse`
func (s *Server) SignTransaction(response http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	var signReq SignTransactionReq
	err := decoder.Decode(&signReq)
	if err != nil {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			err.Error(),
		})
	}

	// we assume that the data to sign in base64 encoded to fit in a JSON doc,
	// so we dcode it here
	decoded, err := base64.StdEncoding.DecodeString(signReq.Data)
	if err != nil {
		WriteErrorResponse(response, http.StatusBadRequest, []string{
			err.Error(),
		})
	}
	ctx := request.Context()
	signature, err := s.singin.SignTransaction(ctx, signReq.ID, decoded)
	if err != nil {
		slog.Error("SignTransaction", "error", fmt.Sprintf("%+v", err))
		WriteInternalError(response)
	}
	WriteAPIResponse(response, http.StatusCreated, signature)
}

func (s *Server) ListTransactions(response http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	list, err := s.singin.ListSignatures(ctx)
	if err != nil {
		slog.Error("ListTransactions", "error", fmt.Sprintf("%+v", err))
		WriteInternalError(response)
	}
	WriteAPIResponse(response, http.StatusOK, SigListToHTTP(list))
}

type APISignature struct {
	Value           string    `json:"value"`
	DeviceID        string    `json:"device_id"`
	SignatureNumber int       `json:"signature_number"`
	SignedData      string    `json:"signed_data"`
	Timestamp       time.Time `json:"ts"`
}

func SigListToHTTP(list []domain.Signature) []APISignature {
	ret := make([]APISignature, 0, len(list))
	for _, sig := range list {
		ret = append(ret, APISignature{
			Value:           sig.Value,
			DeviceID:        sig.DeviceID,
			SignatureNumber: sig.SignatureNumber,
			SignedData:      sig.SignedData,
			Timestamp:       sig.Timestamp,
		})
	}
	return ret

}
