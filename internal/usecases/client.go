package usecases

import (
	"errors"

	"loadbalancer/internal/domain"
	"loadbalancer/internal/interfaces/repositories"
)

type ClientManager struct {
	repo repositories.ClientRepository
}

func NewClientManager(repo repositories.ClientRepository) *ClientManager {
	return &ClientManager{repo: repo}
}

func (m *ClientManager) RegisterClient(id string, capacity, ratePerSec int) (*domain.Client, error) {
	if id == "" {
		return nil, errors.New("client ID cannot be empty")
	}
	if capacity <= 0 {
		return nil, errors.New("capacity must be positive")
	}
	if ratePerSec <= 0 {
		return nil, errors.New("rate must be positive")
	}

	client := domain.NewClient(id, capacity, ratePerSec)
	if err := m.repo.Save(client); err != nil {
		return nil, err
	}
	return client, nil
}

func (m *ClientManager) UpdateClient(id string, capacity, ratePerSec int) (*domain.Client, error) {
	client, err := m.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if capacity > 0 {
		client.Capacity = capacity
	}
	if ratePerSec > 0 {
		client.RatePerSec = ratePerSec
	}

	if err := m.repo.Save(client); err != nil {
		return nil, err
	}
	return client, nil
}

func (m *ClientManager) DeleteClient(id string) error {
	return m.repo.Delete(id)
}

func (m *ClientManager) GetClient(id string) (*domain.Client, error) {
	return m.repo.FindByID(id)
}

func (m *ClientManager) ListClients() ([]*domain.Client, error) {
	return m.repo.FindAll()
}