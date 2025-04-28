package httputil

import (
	"net/http"
	"net/url"
	"time"
)

type HealthChecker interface {
    Check(*url.URL) bool  // Изменено на указатель на url.URL
    Stop()
}

// Реализация
type healthChecker struct {
    interval time.Duration
    stopChan chan struct{}
}

func NewHealthChecker(interval time.Duration) HealthChecker {
    return &healthChecker{
        interval: interval,
        stopChan: make(chan struct{}),
    }
}

// Check выполняет проверку здоровья сервера
func (h *healthChecker) Check(u *url.URL) bool {
	client := http.Client{Timeout: http.DefaultClient.Timeout}
	resp, err := client.Get(u.String() + "/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func (h *healthChecker) Stop() {
    close(h.stopChan)
}