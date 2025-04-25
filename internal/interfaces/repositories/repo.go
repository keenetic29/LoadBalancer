package repositories

import "loadbalancer/internal/domain"

type ServerRepository interface {
	GetAll() ([]*domain.Server, error)
	GetNext() (*domain.Server, error)
	MarkUnhealthy(server *domain.Server) error
}