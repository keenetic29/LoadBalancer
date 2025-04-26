package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"loadbalancer/internal/config"
	"loadbalancer/internal/handlers"
	"loadbalancer/internal/repositories"
	"loadbalancer/internal/usecases"
	"loadbalancer/pkg/httputil"
)

// Запускает тестовый сервер на указанном порту
func startTestServer(port int) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "OK from %d", port)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Response from backend server on port %d", port)
	}

	addr := fmt.Sprintf(":%d", port)
	log.Printf("Starting test server on port %d", port)
	go func() {
		if err := http.ListenAndServe(addr, http.HandlerFunc(handler)); err != nil {
			log.Fatalf("Test server on port %d failed: %v", port, err)
		}
	}()
}


func main() {

	// Запускаем тестовые серверы
	ports := []int{8081, 8082, 8083}
	var wg sync.WaitGroup
	
	for _, port := range ports {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			startTestServer(p)
		}(port)
	}
	
	wg.Wait()

	// Загрузка конфигурации
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Инициализация зависимостей
	serverRepo := repositories.NewMemoryServerRepository(cfg.Backends)
	healthChecker := httputil.NewHealthChecker(2 * time.Second)
	lbUseCase := usecases.NewLoadBalancer(serverRepo, healthChecker)
	handler := handlers.NewLoadBalancerHandler(lbUseCase)

	// Запуск сервера
	log.Printf("Starting load balancer on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}