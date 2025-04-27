package ratelimiter

import (
	"net"
	"sync"
	"time"
)

// LimiterManager управляет токен-бакетами для разных клиентов.
type LimiterManager struct {
	buckets map[string]*TokenBucket
	mu      sync.Mutex

	// Настройки по умолчанию
	defaultCapacity   int
	defaultRefillRate int
	refillPeriod      time.Duration
}

// NewLimiterManager создает новый менеджер лимитеров.
func NewLimiterManager(capacity, refillRate int, refillPeriod time.Duration) *LimiterManager {
	return &LimiterManager{
		buckets:           make(map[string]*TokenBucket),
		defaultCapacity:   capacity,
		defaultRefillRate: refillRate,
		refillPeriod:      refillPeriod,
	}
}

// getOrCreateBucket возвращает бакет для клиента, создавая его при необходимости.
func (m *LimiterManager) getOrCreateBucket(clientID string) *TokenBucket {
	m.mu.Lock()
	defer m.mu.Unlock()

	bucket, exists := m.buckets[clientID]
	if !exists {
		bucket = NewTokenBucket(m.defaultCapacity, m.defaultRefillRate, m.refillPeriod)
		m.buckets[clientID] = bucket
	}
	return bucket
}

// Allow проверяет, можно ли пропустить запрос от клиента.
func (m *LimiterManager) Allow(ip net.IP) bool {
	clientID := ip.String()
	bucket := m.getOrCreateBucket(clientID)
	return bucket.Allow()
}
