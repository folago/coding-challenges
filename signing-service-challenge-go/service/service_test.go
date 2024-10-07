package service_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService(t *testing.T) {
	ctx := context.Background()
	t.Run("ConcurrentSign", func(t *testing.T) {
		var (
			wg     sync.WaitGroup
			n      = 10000
			data   = []byte("some data to sign")
			before = time.Now()
			devID  = "aaaa-bbbb-cccc-dddd"
		)
		type sigtime struct {
			count int
			ts    time.Time
		}

		svc := service.NewSignerService(
			persistence.NewMemDeviceRepository(),
			persistence.NewMemSignatureRepository())

		err := svc.CreateSigner(ctx, devID, "RSA", "device")
		require.NoError(t, err)

		// here we store the signatures
		results := make([]domain.Signature, n)
		for range n {
			wg.Add(1)
			go func() {
				defer wg.Done()
				sign, err := svc.SignTransaction(ctx, devID, data)
				require.NoError(t, err)
				// the Signature number starts at 1, so to have less fiddling with
				// indices later we make it start a 0
				sign.SignatureNumber -= 1
				results[sign.SignatureNumber] = sign
			}()
		}
		wg.Wait()
		// now we test for monotonicity
		for i, sig := range results {
			assert.Equal(t, i, sig.SignatureNumber)
			assert.Truef(t, before.Before(sig.Timestamp), "%v is not before %v", before, sig.Timestamp)
			sig.Timestamp = before
		}
	})
}
