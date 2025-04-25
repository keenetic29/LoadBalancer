package repositories

import (
	"sync"
	
	"loadbalancer/internal/domain"
)

type MemoryServerRepository struct {
	servers []*domain.Server
	current uint64
	mu      sync.Mutex
}

func NewMemoryServerRepository(backends []string) *MemoryServerRepository {
	var servers []*domain.Server
	for _, backend := range backends {
		if server, err := domain.NewServer(backend); err == nil {
			servers = append(servers, server)
		}
	}
	return &MemoryServerRepository{servers: servers}
}

func (r *MemoryServerRepository) GetAll() ([]*domain.Server, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.servers, nil
}

func (r *MemoryServerRepository) GetNext() (*domain.Server, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.servers) == 0 {
		return nil, nil
	}

	originalIndex := r.current
	for {
		server := r.servers[r.current%uint64(len(r.servers))]
		r.current++

		if server.Healthy {
			return server, nil
		}

		if r.current%uint64(len(r.servers)) == originalIndex%uint64(len(r.servers)) {
			return nil, nil
		}
	}
}

func (r *MemoryServerRepository) MarkUnhealthy(server *domain.Server) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, s := range r.servers {
		if s.URL.String() == server.URL.String() {
			s.Healthy = false
			break
		}
	}
	return nil
}