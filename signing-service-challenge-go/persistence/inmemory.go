package persistence

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
)

// memDevRepo is just a map[uuid]Device, the uuid is used in its
// string form
type memDevRepo struct {
	devices map[string]*domain.Device
	mu      sync.Mutex
}

func NewMemDeviceRepository() *memDevRepo {
	return &memDevRepo{
		devices: map[string]*domain.Device{},
	}
}

func (r *memDevRepo) CreateDevice(ctx context.Context, device *domain.Device) error {
	foundDevice, err := r.FindByID(ctx, device.ID)
	if err != nil && !IsNotFoundErr(err) {
		return fmt.Errorf("error creating device: %w", err)
	}
	if err == nil && foundDevice != nil {
		return fmt.Errorf("error creating device: another device with ID: %q already exists", foundDevice.ID)
	}
	err = r.Save(ctx, device)
	if err != nil {
		return fmt.Errorf("error saving device: %w", err)
	}

	return nil
}

func (r *memDevRepo) FindByID(ctx context.Context, id string) (*domain.Device, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	d, ok := r.devices[id]
	if !ok {
		return nil, fmt.Errorf("entity not found")
	}
	return d, nil
}

func (r *memDevRepo) FindByLabel(ctx context.Context, label string) (*domain.Device, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, d := range r.devices {
		if d.Label == label {
			return d, nil
		}
	}
	return nil, fmt.Errorf("entity not found")
}

func (r *memDevRepo) Save(ctx context.Context, dev *domain.Device) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.devices[dev.ID] = dev
	return nil
}

func (r *memDevRepo) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.devices, id)

	return nil
}

func (r *memDevRepo) List(ctx context.Context) ([]domain.Device, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	ret := []domain.Device{}
	for _, v := range r.devices {
		ret = append(ret, *v)
	}

	return ret, nil
}

// memSignRepo is just a map[deviceID]signature
type memSignRepo struct {
	signatures map[string]domain.Signature
	mu         sync.Mutex
}

func NewMemSignatureRepository() *memSignRepo {
	return &memSignRepo{
		signatures: map[string]domain.Signature{},
	}
}

func (r *memSignRepo) FindByDeviceID(ctx context.Context, device string) (*domain.Signature, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, s := range r.signatures {
		if s.DeviceID == device {
			return &s, nil
		}
	}
	return nil, fmt.Errorf("entity not found")
}

func (r *memSignRepo) FindByValue(ctx context.Context, signature string) (*domain.Signature, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	s, ok := r.signatures[signature]
	if !ok {
		return nil, fmt.Errorf("entity not found")
	}
	return &s, fmt.Errorf("entity not found")
}

func (r *memSignRepo) Save(ctx context.Context, signature domain.Signature) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.signatures[signature.Value] = signature

	return nil
}

func (r *memSignRepo) List(ctx context.Context) ([]domain.Signature, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	ret := []domain.Signature{}
	for _, v := range r.signatures {
		ret = append(ret, v)
	}

	return ret, nil
}

func IsNotFoundErr(err error) bool {
	if err == nil {
		return false
	}
	if strings.Contains(err.Error(), "not found") {
		return true
	}
	return false
}
