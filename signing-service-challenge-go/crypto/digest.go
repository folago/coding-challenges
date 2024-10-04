package crypto

import (
	"crypto"
	_ "golang.org/x/crypto/sha3"
)

// DigestAlgorithm is the algorithm used to hash the payload before signing.
// The SHA3 family seems to have the blessing of NIST:
// https://csrc.nist.gov/projects/hash-functions
// we picked the 256 version because (I like powers of 2 :P) the most common
// elliptic curves should have a bit-length of the private key's curve order of
// at least 256, and the length of the digest to sign should be <= of that or
// it gets truncated in ECDSA. https://pkg.go.dev/crypto/ecdsa@go1.21.0#SignASN1
const DigestAlgorithm = crypto.SHA3_256

// Digest hashes the data with HashAlgo
func Digest(data []byte) []byte {
	hash := DigestAlgorithm.New()
	hash.Write(data)
	return hash.Sum(nil)
}
