package service

import (
	"context"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

type Signer interface {
	CreateSigner(ctx context.Context, id, algorithm, label string) error
	ListDevices(ctx context.Context) ([]domain.Device, error)
	DeleteSigner(ctx context.Context, id string) error

	SignTransaction(ctx context.Context, deviceId string, data []byte) (domain.Signature, error)
	ListSignatures(ctx context.Context) ([]domain.Signature, error)
}
