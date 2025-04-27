package handlers

import (
	"net"
	"net/http"
	"time"

	"loadbalancer/internal/interfaces/handlers"
	"loadbalancer/internal/interfaces/repositories"
	"loadbalancer/internal/interfaces/usecases"
	"loadbalancer/internal/ratelimiter"
)

type loadBalancerHandler struct {
    useCase usecases.LoadBalancerUseCase
    limiterManager *ratelimiter.LimiterManager
}

func NewLoadBalancerHandler(
	uc usecases.LoadBalancerUseCase,
	clientRepo repositories.ClientRepository,
	defaultCapacity int,
	defaultRatePerSec int,
	refillPeriod time.Duration,
) handlers.LoadBalancerHandler {
	limiter := ratelimiter.NewLimiterManager(
		clientRepo,
		defaultCapacity,
		defaultRatePerSec,
		refillPeriod,
	)
	return &loadBalancerHandler{
		useCase:        uc,
		limiterManager: limiter,
	}
}
func (h *loadBalancerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, "Failed to parse IP address", http.StatusBadRequest)
		return
	}

	clientIP := net.ParseIP(ip)
	if clientIP == nil {
		http.Error(w, "Invalid IP address", http.StatusBadRequest)
		return
	}

	if !h.limiterManager.Allow(clientIP) {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	// Иначе — пропускаем запрос
    h.useCase.HandleRequest(w, r)
}