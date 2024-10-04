package persistence

import (
	"context"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

// DeviceRepository is the interface implemented by types that can store and
// find Devices.
type DeviceRepository interface {
	CreateDevice(ctx context.Context, dev *domain.Device) error
	FindByID(ctx context.Context, id string) (*domain.Device, error)
	// NOTE: is label unique?
	FindByLabel(ctx context.Context, label string) (*domain.Device, error)
	Save(ctx context.Context, dev *domain.Device) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.Device, error)
}

// SignatureRepository is the intereface implemented by types that can store
// ad find signatures.
type SignatureRepository interface {
	FindByDeviceID(ctx context.Context, device string) (*domain.Signature, error)
	FindByValue(ctx context.Context, signature string) (*domain.Signature, error)
	Save(ctx context.Context, signature domain.Signature) error
	List(ctx context.Context) ([]domain.Signature, error)
}
