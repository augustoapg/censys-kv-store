package store

import (
	"errors"
	"time"
)

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
	GetKvByKey(key string) (KV, error)
	DeleteKvByKey(key string) error
}

var ErrKeyNotFound = errors.New("key not found")

// GetKvByKey retrieves a key-value pair from the store.
// If key is not found, or is soft deleted, return an error.
// Returns the key-value pair.
func (s *InMemoryKVStore) GetKvByKey(key string) (KV, error) {
	kv, ok := s.store[key]

	if !ok || kv.DeletedAt != nil {
		return KV{}, ErrKeyNotFound
	}

	return kv, nil
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

// DeleteKvByKey soft deletes a key-value pair from the store.
// Returns an error if the key is not found or already soft deleted.
func (s *InMemoryKVStore) DeleteKvByKey(key string) error {
	kv, err := s.GetKvByKey(key)

	if err != nil {
		return err
	}

	now := time.Now()
	kv.DeletedAt = &now
	kv.UpdatedAt = now

	s.store[key] = kv

	return nil
}
