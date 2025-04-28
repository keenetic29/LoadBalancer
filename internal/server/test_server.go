// loadbalancer/internal/server/test_server.go
package server

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

type TestServer struct {
	port     int
	server   *http.Server
	isAlive  bool
	mu       sync.RWMutex
}

func NewTestServer(port int) *TestServer {
	s := &TestServer{
		port: port,
		isAlive: true,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleRequest)
	mux.HandleFunc("/health", s.handleHealthCheck)

	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go s.simulateFailures()
	return s
}

func (s *TestServer) Start(ready chan<- struct{}) error {
	go func() {
		log.Printf("Starting test server on port %d", s.port)
		ready <- struct{}{}
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Test server on port %d failed: %v", s.port, err)
		}
	}()
	return nil
}

func (s *TestServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if !s.isAlive {
		http.Error(w, "Server temporarily down", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Response from backend server on port %d", s.port)
}

func (s *TestServer) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	if !s.isAlive {
		http.Error(w, "Server is down", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "OK from %d", s.port)
}

func (s *TestServer) simulateFailures() {
	for {
		time.Sleep(time.Duration(5+rand.Intn(10)) * time.Second)

		s.mu.Lock()
		s.isAlive = !s.isAlive
		state := "DOWN"
		if s.isAlive {
			state = "UP"
		}
		log.Printf("[Server %d] Now %s", s.port, state)
		s.mu.Unlock()
	}
}