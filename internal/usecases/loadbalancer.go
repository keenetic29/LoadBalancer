package usecases

import (
	"log"
	"net/http"
	"net/http/httputil"

	
	"loadbalancer/internal/interfaces/repositories"
	util "loadbalancer/pkg/httputil"
)

type LoadBalancer struct {
	serverRepo    repositories.ServerRepository
	healthChecker util.HealthChecker  
	proxy         *httputil.ReverseProxy
}

func NewLoadBalancer(repo repositories.ServerRepository, checker util.HealthChecker) *LoadBalancer {
	return &LoadBalancer{
		serverRepo:    repo,
		healthChecker: checker,
		proxy:         &httputil.ReverseProxy{},
	}
}

func (lb *LoadBalancer) HandleRequest(w http.ResponseWriter, r *http.Request) {
	server, err := lb.serverRepo.GetNext()
	if err != nil || server == nil {
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}

	if !lb.healthChecker.Check(server.URL) {
		lb.serverRepo.MarkUnhealthy(server)
		lb.HandleRequest(w, r) // Попробовать другой сервер
		return
	}

	lb.proxy.Director = func(req *http.Request) {
		req.URL.Scheme = server.URL.Scheme
		req.URL.Host = server.URL.Host
		req.Host = server.URL.Host
	}

	log.Printf("Proxying request to %s", server.URL.String())
	lb.proxy.ServeHTTP(w, r)
}