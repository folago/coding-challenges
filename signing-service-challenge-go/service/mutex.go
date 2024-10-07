package service

import "sync"

// DeviceMutex is a map of locks to guarantee mutual exclusive access to the
// devices at a service level, a sync.Map might be a better choice in the long
// run but this has a smaller API surface
type DeviceMutex struct {
	devices map[string]*sync.Mutex
	mu      *sync.RWMutex
}

func NewDeviceMutex() *DeviceMutex {
	return &DeviceMutex{
		devices: map[string]*sync.Mutex{},
		mu:      &sync.RWMutex{},
	}
}

func (dm DeviceMutex) Get(id string) (*sync.Mutex, bool) {
	dm.mu.RLock()
	defer dm.mu.RUnlock()

	devMuex, ok := dm.devices[id]

	return devMuex, ok
}

func (dm DeviceMutex) Put(id string, devMu *sync.Mutex) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	dm.devices[id] = devMu
}

func (dm DeviceMutex) Delete(id string) {
	dm.mu.Lock()
	defer dm.mu.Unlock()

	delete(dm.devices, id)
}
