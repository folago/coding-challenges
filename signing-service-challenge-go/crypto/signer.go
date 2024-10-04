package crypto

//go:generate go run github.com/dmarkham/enumer -type=SignAlgorithm -json

// SignAlgorithm is the enum for the signing algorithms supported
// NOTE: to add a new algorithm add a constant and run `go generate`, in case
// it is needed also run `go get github.com/dmarkham/enumer` before
type SignAlgorithm int8

const (
	RSA SignAlgorithm = iota
	ECC
)

// Signer defines a contract for different types of signing implementations.
type Signer interface {
	Sign(dataToBeSigned []byte) ([]byte, error)
}
