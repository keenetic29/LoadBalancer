package usecases

import (
	"net/http"
)

type LoadBalancerUseCase interface {
	HandleRequest(w http.ResponseWriter, r *http.Request)
}