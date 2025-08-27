package store

import "time"

type KV struct {
	Key       string     `json:"key"`
	Value     string     `json:"value"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type InMemoryKVStore struct {
	store map[string]KV
}

func NewInMemoryKVStore() KVStore {
	return &InMemoryKVStore{
		store: make(map[string]KV),
	}
}

type KVStore interface {
	UpsertKv(kv *KV) (KV, error)
}

// UpsertKv creates or updates a key-value pair in the store.
// If the key is soft deleted, it will be restored.
// Returns the created or updated key-value pair.
func (s *InMemoryKVStore) UpsertKv(kv *KV) (KV, error) {
	key := kv.Key
	value := kv.Value

	createdAt := time.Now()

	if _, ok := s.store[key]; ok {
		createdAt = s.store[key].CreatedAt
	}

	s.store[key] = KV{
		Key:       key,
		Value:     value,
		CreatedAt: createdAt,
		UpdatedAt: time.Now(),
		DeletedAt: nil,
	}

	return s.store[key], nil
}
