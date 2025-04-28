package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"loadbalancer/internal/config"
	
	"loadbalancer/internal/server"
)

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

func startTestServers(ports []int) {
	ready := make(chan struct{})
	for _, port := range ports {
		ts := server.NewTestServer(port)
		go ts.Start(ready)
	}

	// Ждем готовности всех серверов
	for i := 0; i < len(ports); i++ {
		<-ready
	}
	log.Println("All test servers are ready!")
}



func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Извлечение портов бэкендов
	ports, err := extractPortsFromBackends(cfg.Backends)
	if err != nil {
		log.Fatalf("Failed to extract ports from backends: %v", err)
	}

	// Запуск тестовых серверов
	startTestServers(ports)

	// Создание и запуск сервера балансировщика
	lbServer := server.NewLoadBalancerServer(cfg)
	if err := lbServer.Start(); err != nil {
		log.Fatalf("Failed to start load balancer server: %v", err)
	}
	log.Printf("Load balancer started on port %s", cfg.Port)

	// Обработка сигналов завершения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down server...")
	if err := lbServer.Stop(); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}
	log.Println("Server gracefully stopped")
}