package crypto_test

import (
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/stretchr/testify/require"
)

// TODO: maybe this check should be an init()?
func TestHashAlgorithmAvailable(t *testing.T) {
	require.True(t, crypto.DigestAlgorithm.Available())
}
