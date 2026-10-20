package store

import (
	"errors"
	"sync"
)

type Monitor struct {
	ID              string `json:"id"`
	Url             string `json:"url"`
	IntervalSeconds int    `json:"interval_seconds"`
	Enabled         bool   `json:"enabled"`
}

type Store struct {
	mu   sync.RWMutex
	data map[string]Monitor
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]Monitor),
	}
}

func (s *Store) Set(key string, value Monitor) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.data[key]
	if exists {
		return false, errors.New("Key is already exists")
	}

	s.data[key] = value

	return true, nil
}

func (s *Store) Get(key string) (Monitor, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, exists := s.data[key]
	if !exists {
		return Monitor{}, errors.New("monitor does not exist")
	}

	return val, nil
}

func (s *Store) Update(key string, value Monitor) (Monitor, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.data[key]; !exists {
		return Monitor{}, errors.New("monitor does not exist")
	}

	value.ID = key
	s.data[key] = value

	return value, nil
}

func (s *Store) Delete(key string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exists := s.data[key]
	if !exists {
		return false, errors.New("Key does not exists")
	}

	delete(s.data, key)

	return true, nil
}

func (s *Store) List() []Monitor {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Monitor, 0, len(s.data))

	for _, val := range s.data {
		result = append(result, val)
	}

	return result
}
