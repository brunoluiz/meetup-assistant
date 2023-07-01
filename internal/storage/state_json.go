package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"golang.org/x/exp/slog"
)

type StateJSON struct {
	path  string
	state map[string]string
	mu    sync.Mutex
}

func NewStateJSON(path string) (*StateJSON, error) {
	obj, err := os.ReadFile(path)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		obj = []byte("{}")
	}

	state := map[string]string{}
	if err := json.Unmarshal(obj, &state); err != nil {
		return nil, err
	}

	return &StateJSON{
		path:  path,
		state: state,
		mu:    sync.Mutex{},
	}, nil
}

func (m *StateJSON) key(audience, key string) string {
	return fmt.Sprintf("%s:%s", audience, key)
}

func (m *StateJSON) Get(ctx context.Context, audience, key string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	slog.DebugCtx(ctx, "StateJSON.Get", "audience", audience, "key", key)
	return m.state[m.key(audience, key)], nil
}

func (m *StateJSON) Save(ctx context.Context, audience, key, state string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.state[m.key(audience, key)] = state

	obj, err := json.Marshal(m.state)
	if err != nil {
		return err
	}

	slog.DebugCtx(ctx, "StateJSON.Save", "audience", audience, "key", key, "state", state)
	return os.WriteFile(m.path, obj, 0644)
}
