package repositories

import "loadbalancer/internal/domain"

type ServerRepository interface {
	GetNext() (*domain.Server, error)
	MarkUnhealthy(server *domain.Server) error
	Count() int
	GetAll() []*domain.Server
	UpdateHealth(server *domain.Server, healthy bool)
}
