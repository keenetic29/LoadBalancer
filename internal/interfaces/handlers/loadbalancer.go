package handlers

import (
	"net/http"
	
	"loadbalancer/internal/interfaces/usecases"
)

type LoadBalancerHandler struct {
	useCase usecases.LoadBalancerUseCase
}

func NewLoadBalancerHandler(uc usecases.LoadBalancerUseCase) *LoadBalancerHandler {
	return &LoadBalancerHandler{useCase: uc}
}

func (h *LoadBalancerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.useCase.HandleRequest(w, r)
}