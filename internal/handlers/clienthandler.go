package handlers

import (
	"encoding/json"
	"net/http"

	"loadbalancer/internal/interfaces/usecases"
)

type ClientHandler struct {
	useCase usecases.ClientUseCase
}

func NewClientHandler(uc usecases.ClientUseCase) *ClientHandler {
	return &ClientHandler{useCase: uc}
}

type clientRequest struct {
	ID         string `json:"client_id"`
	Capacity   int    `json:"capacity"`
	RatePerSec int    `json:"rate_per_sec"`
}

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (h *ClientHandler) RegisterClient(w http.ResponseWriter, r *http.Request) {
	var req clientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	client, err := h.useCase.RegisterClient(req.ID, req.Capacity, req.RatePerSec)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, client)
}

func (h *ClientHandler) UpdateClient(w http.ResponseWriter, r *http.Request) {
	var req clientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	client, err := h.useCase.UpdateClient(req.ID, req.Capacity, req.RatePerSec)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, client)
}

func (h *ClientHandler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "client ID is required")
		return
	}

	if err := h.useCase.DeleteClient(id); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ClientHandler) GetClient(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		respondWithError(w, http.StatusBadRequest, "client ID is required")
		return
	}

	client, err := h.useCase.GetClient(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, client)
}

func (h *ClientHandler) ListClients(w http.ResponseWriter, r *http.Request) {
	clients, err := h.useCase.ListClients()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, clients)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errorResponse{
		Code:    code,
		Message: message,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}