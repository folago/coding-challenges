package domain

import "time"

// Signature holds the data signed, and the signature, and the ID of the device
// that signed the data
type Signature struct {
	Value           string
	DeviceID        string
	SignatureNumber int
	SignedData      string
	Timestamp       time.Time
}
