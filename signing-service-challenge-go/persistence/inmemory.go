package persistence

import (
	"context"
	"sync"

	"github.com/gdelafuente/fiskaly-coding-challenges/signing-service-challenge-go/domain"
)

type InMemoryDeviceRepository struct {
	data map[string]domain.Device
	lock sync.RWMutex
}

func NewInMemoryDeviceRepository() *InMemoryDeviceRepository {
	return &InMemoryDeviceRepository{
		data: make(map[string]domain.Device),
	}
}

func (r *InMemoryDeviceRepository) Save(ctx context.Context, device domain.Device) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.data[device.ID()] = device
	return nil
}

func (r *InMemoryDeviceRepository) Update(ctx context.Context, device domain.Device, expectedVersion int) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	existingDevice, ok := r.data[device.ID()]
	if !ok {
		return domain.ErrDeviceNotFound
	}
	if existingDevice.Version() != expectedVersion {
		return domain.ErrDeviceVersionMismatch
	}

	r.data[device.ID()] = device
	return nil
}

func (r *InMemoryDeviceRepository) FindByID(ctx context.Context, id string) (domain.Device, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	device, ok := r.data[id]
	if !ok {
		return domain.Device{}, domain.ErrDeviceNotFound
	}
	return device, nil
}

func (r *InMemoryDeviceRepository) ListAll(ctx context.Context) ([]domain.Device, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	// TODO: avoid iterating over all the values
	result := make([]domain.Device, 0, len(r.data))
	for _, device := range r.data {
		result = append(result, device)
	}
	return result, nil
}
