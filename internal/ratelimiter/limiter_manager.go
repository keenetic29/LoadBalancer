package ratelimiter

import (
	"net"
	"sync"
	"time"

	"loadbalancer/internal/interfaces/repositories"
)

type LimiterManager struct {
	buckets        map[string]*TokenBucket
	mu             sync.Mutex
	clientRepo     repositories.ClientRepository
	defaultLimiter *TokenBucket
}

func NewLimiterManager(clientRepo repositories.ClientRepository, defaultCapacity, defaultRefillRate int, refillPeriod time.Duration) *LimiterManager {
	return &LimiterManager{
		buckets:        make(map[string]*TokenBucket),
		clientRepo:     clientRepo,
		defaultLimiter: NewTokenBucket(defaultCapacity, defaultRefillRate, refillPeriod),
	}
}

func (m *LimiterManager) getOrCreateBucket(clientID string) (*TokenBucket, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if bucket, exists := m.buckets[clientID]; exists {
		return bucket, nil
	}

	// попытка найти лимиты
	client, err := m.clientRepo.FindByID(clientID)
	if err == nil {
		bucket := NewTokenBucket(client.Capacity, client.RatePerSec, client.RefillPeriod)
		m.buckets[clientID] = bucket
		return bucket, nil
	}

	
	return m.defaultLimiter, nil
}

func (m *LimiterManager) Allow(ip net.IP) bool {
	clientID := ip.String()
	bucket, err := m.getOrCreateBucket(clientID)
	if err != nil {
		return false
	}
	return bucket.Allow()
}