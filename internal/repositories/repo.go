package repositories

import (
	"fmt"
	"sync"

	"loadbalancer/internal/domain"
)

type MemoryServerRepository struct {
	servers []*domain.Server
	current int
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

func (r *MemoryServerRepository) GetAll() []*domain.Server {
	r.mu.Lock()
	defer r.mu.Unlock()

	servers := make([]*domain.Server, len(r.servers))
	copy(servers, r.servers)
	return servers
}

func (r *MemoryServerRepository) UpdateHealth(server *domain.Server, healthy bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, s := range r.servers {
		if s.URL.String() == server.URL.String() {
			s.Healthy = healthy
			break
		}
	}
}


func (r *MemoryServerRepository) GetNext() (*domain.Server, error) {
    r.mu.Lock()
	defer r.mu.Unlock()

	if len(r.servers) == 0 {
		return nil, fmt.Errorf("no servers available")
	}

	startIdx := r.current
	for {
		server := r.servers[r.current]
		r.current = (r.current + 1) % len(r.servers) // Круговая очередь

		if server.Healthy {
			return server, nil
		}

		// Прошли все серверы и не нашли здоровый
		if r.current == startIdx {
			return nil, fmt.Errorf("no healthy servers available")
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

func (r *MemoryServerRepository) Count() int {
    r.mu.Lock()
    defer r.mu.Unlock()
    return len(r.servers)
}
