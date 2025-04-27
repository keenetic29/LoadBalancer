package repositories

import "loadbalancer/internal/domain"

type ClientRepository interface {
	Save(client *domain.Client) error
	FindByID(id string) (*domain.Client, error)
	Delete(id string) error
	FindAll() ([]*domain.Client, error)
}