package usecases

import "loadbalancer/internal/domain"

type ClientUseCase interface {
	RegisterClient(id string, capacity, ratePerSec int) (*domain.Client, error)
	UpdateClient(id string, capacity, ratePerSec int) (*domain.Client, error)
	DeleteClient(id string) error
	GetClient(id string) (*domain.Client, error)
	ListClients() ([]*domain.Client, error)
}