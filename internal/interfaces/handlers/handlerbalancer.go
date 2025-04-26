package handlers

import "net/http"

type LoadBalancerHandler interface {
    ServeHTTP(w http.ResponseWriter, r *http.Request)
}