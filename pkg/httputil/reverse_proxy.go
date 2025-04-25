package httputil

import (
	"net/http"
	"net/url"
	"time"
)

// HealthChecker определяет интерфейс для проверки здоровья
type HealthChecker interface {
	Check(*url.URL) bool
}

// httpHealthChecker реализует HealthChecker
type httpHealthChecker struct {
	timeout time.Duration
}

// NewHealthChecker создает новый экземпляр HealthChecker
func NewHealthChecker(timeout time.Duration) HealthChecker {
	return &httpHealthChecker{timeout: timeout}
}

// Check выполняет проверку здоровья сервера
func (h *httpHealthChecker) Check(u *url.URL) bool {
	client := http.Client{Timeout: h.timeout}
	resp, err := client.Get(u.String() + "/health")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

