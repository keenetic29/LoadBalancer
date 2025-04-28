// loadbalancer/internal/server/server.go
package server

import (
	"context"
	"net/http"
	"sync"
	"time"

	"loadbalancer/internal/config"
	"loadbalancer/internal/handlers"
	"loadbalancer/internal/repositories"
	"loadbalancer/internal/usecases"
	util "loadbalancer/pkg/httputil"
)

type LoadBalancerServer struct {
	server     		*http.Server
	healthChecker 	util.HealthChecker
	wg         		sync.WaitGroup
}

func NewLoadBalancerServer(cfg *config.Config) *LoadBalancerServer {
	// Инициализация зависимостей
	serverRepo := repositories.NewMemoryServerRepository(cfg.Backends)
	clientRepo := repositories.NewMemoryClientRepository(cfg.ClientsDB)
	healthChecker := util.NewHealthChecker(2 * time.Second)

	// Инициализация use cases
	lbUseCase := usecases.NewLoadBalancer(serverRepo, healthChecker)
	clientUseCase := usecases.NewClientManager(clientRepo)
	
	// Инициализация обработчиков
	lbHandler := handlers.NewLoadBalancerHandler(
		lbUseCase,
		clientRepo,
		cfg.RateLimit.DefaultCapacity,
		cfg.RateLimit.DefaultRatePerSec,
		time.Duration(cfg.RateLimit.RefillPeriod)*time.Nanosecond,
	)
	clientHandler := handlers.NewClientHandler(clientUseCase)

	// Настройка маршрутизатора
	mux := http.NewServeMux()
	mux.Handle("/", lbHandler)
	mux.HandleFunc("/clients/register", clientHandler.RegisterClient)
	mux.HandleFunc("/clients/update", clientHandler.UpdateClient)
	mux.HandleFunc("/clients/delete", clientHandler.DeleteClient)
	mux.HandleFunc("/clients/get", clientHandler.GetClient)
	mux.HandleFunc("/clients/list", clientHandler.ListClients)

	return &LoadBalancerServer{
		server: &http.Server{
			Addr:    ":" + cfg.Port,
			Handler: mux,
		},
		healthChecker: healthChecker,
	}
}

func (s *LoadBalancerServer) Start() error {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
	return nil
}

func (s *LoadBalancerServer) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	
	s.healthChecker.Stop()
	s.wg.Wait()
	return nil
}