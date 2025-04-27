package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"

	"loadbalancer/internal/config"
	"loadbalancer/internal/handlers"
	"loadbalancer/internal/repositories"
	"loadbalancer/internal/usecases"
	util "loadbalancer/pkg/httputil"
)

func startTestServer(port int, ready chan<- struct{}) {
	mux := http.NewServeMux()

	isAlive := true
	var mu sync.RWMutex

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mu.RLock()
		defer mu.RUnlock()
		if !isAlive {
			http.Error(w, "Server temporarily down", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Response from backend server on port %d", port)
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		mu.RLock()
		defer mu.RUnlock()
		if !isAlive {
			http.Error(w, "Server is down", http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK from %d", port)
	})

	addr := fmt.Sprintf(":%d", port)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		log.Printf("Starting test server on port %d", port)
		// <- Вот здесь говорим что сервер запущен
		ready <- struct{}{}
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Test server on port %d failed: %v", port, err)
		}
	}()

	// Симуляция падений и восстановления
	go func() {
		for {
			time.Sleep(time.Duration(5+rand.Intn(10)) * time.Second)

			mu.Lock()
			isAlive = !isAlive
			state := "DOWN"
			if isAlive {
				state = "UP"
			}
			log.Printf("[Server %d] Now %s", port, state)
			mu.Unlock()
		}
	}()
}


func extractPortsFromBackends(backends []string) ([]int, error) {
	var ports []int
	for _, backend := range backends {
		parsedURL, err := url.Parse(backend)
		if err != nil {
			return nil, fmt.Errorf("failed to parse backend URL %s: %w", backend, err)
		}

		port := 80 // по умолчанию
		if parsedURL.Port() != "" {
			fmt.Sscanf(parsedURL.Port(), "%d", &port)
		}
		ports = append(ports, port)
	}
	return ports, nil
}



func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Извлекаем порты из backend URL'ов
	ports, err := extractPortsFromBackends(cfg.Backends)
	if err != nil {
		log.Fatalf("Failed to extract ports from backends: %v", err)
	}

	ready := make(chan struct{})

	for _, port := range ports {
		go startTestServer(port, ready)
	}

	// Ждем подтверждения от всех серверов
	for i := 0; i < len(ports); i++ {
		<-ready
	}

	log.Println("All test servers are ready!")


	// Инициализация зависимостей
	serverRepo := repositories.NewMemoryServerRepository(cfg.Backends)
	healthChecker := util.NewHealthChecker(2 * time.Second)
	lbUseCase := usecases.NewLoadBalancer(serverRepo, healthChecker)
	handler := handlers.NewLoadBalancerHandler(lbUseCase)

	// Запуск сервера
	log.Printf("Starting load balancer on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}