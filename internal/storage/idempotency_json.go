package storage

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type IdempotencyJSON struct {
	path  string
	state map[string]bool
	mu    sync.Mutex
}

func NewIdempotencyJSON(path string) (*IdempotencyJSON, error) {
	obj, err := os.ReadFile(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		obj = []byte("{}")
	}

	state := map[string]bool{}
	if err := json.Unmarshal(obj, &state); err != nil {
		return nil, err
	}

	return &IdempotencyJSON{
		path:  path,
		state: state,
		mu:    sync.Mutex{},
	}, nil
}

func (m *IdempotencyJSON) Get(ctx context.Context, key string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.state[key], nil
}

func (m *IdempotencyJSON) Save(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.state[key] = true

	obj, err := json.Marshal(m.state)
	if err != nil {
		return err
	}

	return os.WriteFile(m.path, obj, 0644)
}
