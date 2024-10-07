package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
)

type singerService struct {
	devRepo persistence.DeviceRepository
	sigRepo persistence.SignatureRepository
	devMux  *DeviceMutex
}

func NewSignerService(devRepo persistence.DeviceRepository, sigRepo persistence.SignatureRepository) *singerService {
	return &singerService{
		devRepo: devRepo,
		sigRepo: sigRepo,
		devMux:  NewDeviceMutex(),
	}
}

func (s *singerService) CreateSigner(ctx context.Context, id, algorithm, label string) error {

	dev, err := s.devRepo.FindByID(ctx, id)
	if err != nil && !persistence.IsNotFoundErr(err) {
		return fmt.Errorf("error creating device: %w", err)
	}
	if err == nil && dev != nil {
		return fmt.Errorf("error creating device: another device with ID: %q already exists", id)
	}
	dev, err = domain.NewDevice(id, algorithm, label)
	if err != nil {
		return fmt.Errorf("error creating device: %w", err)
	}
	err = s.devRepo.CreateDevice(ctx, dev)
	if err != nil {
		return fmt.Errorf("error creating device with ID %q: %w", id, err)
	}
	return nil
}
func (s *singerService) DeleteSigner(ctx context.Context, id string) error {
	return s.devRepo.Delete(ctx, id)
}

func (s *singerService) SignTransaction(ctx context.Context, deviceId string, data []byte) (domain.Signature, error) {
	// NOTE: critical section start
	s.lockDevice(deviceId)
	dev, err := s.devRepo.FindByID(ctx, deviceId)
	if err != nil {
		return domain.Signature{}, fmt.Errorf("error signing with device %q: %w", deviceId, err)
	}

	signature, dataSigned, ts, err := dev.Sign(data)
	if err != nil {
		return domain.Signature{}, fmt.Errorf("error signing with device %q: %w", deviceId, err)
	}

	err = s.devRepo.Save(ctx, dev)
	if err != nil {
		return domain.Signature{}, fmt.Errorf("error persisting device %q: %w", deviceId, err)
	}
	count := dev.SignatureCounter
	s.unlockDevice(deviceId)
	// NOTE: critical section end

	ret := domain.Signature{
		Value:           signature,
		DeviceID:        deviceId,
		SignatureNumber: count,
		SignedData:      string(dataSigned),
		Timestamp:       ts,
	}
	err = s.sigRepo.Save(ctx, ret)
	if err != nil {
		return domain.Signature{}, fmt.Errorf("error persisting signature: %w", err)
	}
	return ret, nil
}

func (s *singerService) ListSignatures(ctx context.Context) ([]domain.Signature, error) {
	list, err := s.sigRepo.List(ctx)
	if err != nil {
		return list, fmt.Errorf("error listing sigatures: %w", err)
	}
	return list, nil
}

func (s *singerService) ListDevices(ctx context.Context) ([]domain.Device, error) {
	list, err := s.devRepo.List(ctx)
	if err != nil {
		return list, fmt.Errorf("error listing devices: %w", err)
	}
	return list, nil
}

func (s *singerService) lockDevice(id string) {
	mu, ok := s.devMux.Get(id)
	if !ok {
		currMutex := &sync.Mutex{}
		currMutex.Lock()
		s.devMux.Put(id, currMutex)
		return
	}
	mu.Lock()
}
func (s *singerService) unlockDevice(id string) {
	mu, ok := s.devMux.Get(id)
	if !ok {
		// TODO: need to think on how to handle this situation
		panic("not found device lock while unlocking device")
	}
	mu.Unlock()
	s.devMux.Put(id, mu)
}
