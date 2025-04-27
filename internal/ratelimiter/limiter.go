package ratelimiter
// Реализация Token Bucket

import (
	"sync"
	"time"
)

// TokenBucket представляет отдельный токен-бакет для клиента.
type TokenBucket struct {
	capacity     int
	tokens       int
	refillRate   int           // сколько токенов добавляется за интервал
	refillPeriod time.Duration // интервал пополнения
	lastRefill   time.Time
	mu           sync.Mutex
}

// NewTokenBucket создает новый бакет
func NewTokenBucket(capacity, refillRate int, refillPeriod time.Duration) *TokenBucket {
	return &TokenBucket{
		capacity:     capacity,
		tokens:       capacity,
		refillRate:   refillRate,
		refillPeriod: refillPeriod,
		lastRefill:   time.Now(),
	}
}

// Allow пытается получить токен. Возвращает true, если успешно.
func (b *TokenBucket) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(b.lastRefill)

	// Пополняем токены, если прошло достаточно времени
	if elapsed >= b.refillPeriod {
		refills := int(elapsed / b.refillPeriod)
		b.tokens += refills * b.refillRate
		if b.tokens > b.capacity {
			b.tokens = b.capacity
		}
		b.lastRefill = now
	}

	if b.tokens > 0 {
		b.tokens--
		return true
	}

	return false
}