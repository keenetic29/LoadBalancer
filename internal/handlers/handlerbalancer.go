package handlers

import (
    "net/http"
    
    "loadbalancer/internal/interfaces/handlers"
    "loadbalancer/internal/interfaces/usecases"
)

type loadBalancerHandler struct {
    useCase usecases.LoadBalancerUseCase
}

func NewLoadBalancerHandler(uc usecases.LoadBalancerUseCase) handlers.LoadBalancerHandler {
    return &loadBalancerHandler{useCase: uc}
}

func (h *loadBalancerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    h.useCase.HandleRequest(w, r)
}