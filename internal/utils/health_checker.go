package utils

import (
	"net/http"
	"net/url"
	"time"
)

type HealthChecker interface {
	Check(url *url.URL) bool
}

type HTTPHealthChecker struct {
	timeout time.Duration
}

func NewHTTPHealthChecker(timeout time.Duration) *HTTPHealthChecker {
	return &HTTPHealthChecker{timeout: timeout}
}

func (h *HTTPHealthChecker) Check(u *url.URL) bool {
	resp, err := http.Get(u.String() + "/health")
    if err != nil {
        return false
    }
    defer resp.Body.Close()
    return resp.StatusCode == http.StatusOK
}