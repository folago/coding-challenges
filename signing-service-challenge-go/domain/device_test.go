package domain_test

import (
	"bytes"
	"encoding/base64"
	"strconv"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDevice(t *testing.T) {
	for _, signAlgo := range crypto.SignAlgorithmStrings() {
		t.Run(signAlgo, func(t *testing.T) {
			var (
				dev  *domain.Device
				err  error
				data = []byte("hello world")
			)
			t.Run("Create", func(t *testing.T) {
				dev, err = domain.NewDevice(uuid.New().String(), signAlgo, "coooool device")
				require.NoError(t, err)
				require.Zero(t, dev.SignatureCounter)
				initLastSign := base64.StdEncoding.EncodeToString([]byte(dev.ID))
				require.Equal(t, initLastSign, dev.LastSignature)
			})
			t.Run("Sign", func(t *testing.T) {
				prevSignature := dev.LastSignature
				signature, secData, _, err := dev.Sign(data)
				require.NoError(t, err)

				assert.Equal(t, signature, dev.LastSignature)

				datalist := bytes.Split(secData, []byte{'_'})
				require.Len(t, datalist, 3)

				oldCount, err := strconv.Atoi(string(datalist[0]))
				require.NoError(t, err)
				assert.Zero(t, oldCount)
				assert.Equal(t, 1, dev.SignatureCounter)

				assert.Equal(t, data, datalist[1])
				assert.Equal(t, prevSignature, string(datalist[2]))

			})
		})
	}
	t.Run("FailUnsupportedAlgo", func(t *testing.T) {
		_, err := domain.NewDevice(uuid.New().String(), "DSA", "should not exist device")
		require.Error(t, err)
	})
}
