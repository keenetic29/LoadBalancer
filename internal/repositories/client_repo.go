package repositories

import (
	"encoding/json"
	"errors"
	"os"
	"sync"

	"loadbalancer/internal/domain"
)

type MemoryClientRepository struct {
	clients map[string]*domain.Client
	mu      sync.Mutex
	file    string 
}

func NewMemoryClientRepository(file string) *MemoryClientRepository {
	repo := &MemoryClientRepository{
		clients: make(map[string]*domain.Client),
		file:    file,
	}
	repo.loadFromFile()
	return repo
}

func (r *MemoryClientRepository) Save(client *domain.Client) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.clients[client.ID] = client
	return r.saveToFile()
}

func (r *MemoryClientRepository) FindByID(id string) (*domain.Client, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	client, exists := r.clients[id]
	if !exists {
		return nil, errors.New("client not found")
	}
	return client, nil
}

func (r *MemoryClientRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.clients, id)
	return r.saveToFile()
}

func (r *MemoryClientRepository) FindAll() ([]*domain.Client, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var clients []*domain.Client
	for _, client := range r.clients {
		clients = append(clients, client)
	}
	return clients, nil
}

func (r *MemoryClientRepository) saveToFile() error {
	if r.file == "" {
		return nil
	}

	data, err := json.Marshal(r.clients)
	if err != nil {
		return err
	}
	return os.WriteFile(r.file, data, 0644)
}

func (r *MemoryClientRepository) loadFromFile() error {
	if r.file == "" {
		return nil
	}

	data, err := os.ReadFile(r.file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return json.Unmarshal(data, &r.clients)
}