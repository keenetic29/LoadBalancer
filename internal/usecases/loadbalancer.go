package usecases

import (
	"log"
	"net/http"
	"net/http/httputil"
	"time"

	"loadbalancer/internal/interfaces/repositories"
	util "loadbalancer/pkg/httputil"
)

type LoadBalancer struct {
	serverRepo    repositories.ServerRepository
	healthChecker util.HealthChecker
}

func NewLoadBalancer(repo repositories.ServerRepository, checker util.HealthChecker) *LoadBalancer {
	lb := &LoadBalancer{
		serverRepo:    repo,
		healthChecker: checker,
	}

	// Запуск фоновой проверки здоровья
	go lb.monitorHealth()

	return lb
}

func (lb *LoadBalancer) monitorHealth() {
	ticker := time.NewTicker(3 * time.Second) // Каждые 3 секунды проверяем здоровье
	defer ticker.Stop()

	for {
		<-ticker.C
		lb.checkAllServers()
	}
}

func (lb *LoadBalancer) checkAllServers() {
	servers := lb.serverRepo.GetAll()

	for _, server := range servers {
		isHealthy := lb.healthChecker.Check(server.URL)
		lb.serverRepo.UpdateHealth(server, isHealthy)
	}
}

func (lb *LoadBalancer) HandleRequest(w http.ResponseWriter, r *http.Request) {
	server, err := lb.serverRepo.GetNext()
	if err != nil || server == nil {
		log.Printf("All backend servers are unavailable")
		http.Error(w, "All backend servers are unavailable", http.StatusServiceUnavailable)
		return
	}

	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.URL.Scheme = server.URL.Scheme
			req.URL.Host = server.URL.Host
			req.Host = server.URL.Host
		},
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			log.Printf("Connection to %s failed: %v", server.URL.String(), err)
			lb.serverRepo.MarkUnhealthy(server)
			lb.HandleRequest(w, r) // Пробуем другой сервер
		},
	}

	log.Printf("Proxying request to %s", server.URL.String())
	proxy.ServeHTTP(w, r)
}
