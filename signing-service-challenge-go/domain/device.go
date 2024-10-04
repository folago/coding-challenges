package domain

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
)

type Device struct {
	ID               string
	Algorithm        crypto.SignAlgorithm
	SignatureCounter int // NOTE: is int enough?
	Label            string
	signer           crypto.Signer
	LastSignature    string
}

// Sign will sign the data using the Algorithm provided to the device at
// creation time and return the signature base64 endcoded, the signed data, and an error in
// case something went wrong.
func (d *Device) Sign(data []byte) (string, []byte, time.Time, error) {
	var buff bytes.Buffer
	// TODO: does the counter need 0 padding?
	buff.WriteString(fmt.Sprintf("%d", d.SignatureCounter))
	buff.WriteByte('_')
	buff.Write(data)
	buff.WriteByte('_')
	buff.WriteString(d.LastSignature)

	secureData := buff.Bytes()

	now := time.Now()
	signature, err := d.signer.Sign(secureData)
	if err != nil {
		return "", nil, now, fmt.Errorf("error signing data: %w", err)
	}

	encodedSignature := base64.StdEncoding.EncodeToString(signature)
	d.LastSignature = encodedSignature
	d.SignatureCounter += 1

	return encodedSignature, secureData, now, nil
}

func NewDevice(id, signAlgo, label string) (*Device, error) {
	algo, err := crypto.SignAlgorithmString(signAlgo)
	if err != nil {
		return nil, fmt.Errorf("unsupported signing algorithm %q: %w", signAlgo, err)
	}

	var signer crypto.Signer
	switch algo {
	case crypto.ECC:
		signer, err = crypto.ECCGenerator{}.Generate()
	case crypto.RSA:
		signer, err = crypto.RSAGenerator{}.Generate()
	default:
		// this should not happen but just in case
		return nil, fmt.Errorf("unsupported signing algorithm %q", signAlgo)
	}
	if err != nil {
		return nil, fmt.Errorf("error generating keys for algorithm %q: %w", signAlgo, err)
	}

	ret := &Device{
		ID:            id,
		Algorithm:     algo,
		Label:         label,
		signer:        signer,
		LastSignature: base64.StdEncoding.EncodeToString([]byte(id)),
	}

	return ret, nil
}
