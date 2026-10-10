package storage

import (
	"context"
	"shortener/internal/service"
	"sync"
	"time"
)

type Memory struct {
	mu     sync.RWMutex
	links  map[string]*service.Link
	nextID int64
}

func NewMemory() *Memory {
	return &Memory{
		links:  make(map[string]*service.Link),
		nextID: 1,
	}
}

func (m *Memory) Create(ctx context.Context, link *service.Link) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.links[link.Code]; exists {
		return service.ErrCodeTaken
	}

	link.ID = m.nextID
	m.nextID++
	if link.CreatedAt.IsZero() {
		link.CreatedAt = time.Now()
	}

	stored := *link
	m.links[link.Code] = &stored
	return nil
}

func (m *Memory) GetByCode(ctx context.Context, code string) (*service.Link, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	l, ok := m.links[code]
	if !ok {
		return nil, service.ErrNotFound
	}

	result := *l
	return &result, nil
}
