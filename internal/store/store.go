package store

import (
	"errors"
	"sync"
)

type Store struct {
	mu       sync.RWMutex
	data     map[string]string
	channels map[string]chan map[string]string
}

func New() *Store {
	return &Store{
		data:     make(map[string]string),
		channels: make(map[string]chan map[string]string),
	}
}

func (s *Store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.data[key]
	return v, ok
}

func (s *Store) Del(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}

func (s *Store) CreateChan(name string) (chan map[string]string, error) {
	if name == "" {
		return nil, errors.New("name is empty")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if ch, ok := s.channels[name]; ok {
		return ch, nil
	}

	ch := make(chan map[string]string, 10000)
	s.channels[name] = ch
	return ch, nil
}

func (s *Store) CloseChan(name string) error {
	if name == "" {
		return errors.New("name is empty")
	}

	s.mu.Lock()

	ch, ok := s.channels[name]
	if !ok {
		s.mu.Unlock()
		return errors.New("channel not found")
	}

	delete(s.channels, name)
	s.mu.Unlock()

	close(ch)
	return nil
}

func (s *Store) getOrCreateChan(name string) (chan map[string]string, bool) {

	s.mu.RLock()
	ch, ok := s.channels[name]
	s.mu.RUnlock()
	if ok {
		return ch, false
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if ch, ok := s.channels[name]; ok {
		return ch, false
	}

	ch = make(chan map[string]string, 10000)
	s.channels[name] = ch
	return ch, true
}

func (s *Store) SetToChan(name, key, value string) error {

	if name == "" || key == "" {
		return errors.New("name or key is empty")
	}

	ch, _ := s.getOrCreateChan(name)

	data := map[string]string{key: value}

	select {
	case ch <- data:
		return nil
	default:
		return errors.New("channel is full")
	}
}

func (s *Store) ReadFromChan(name string) (map[string]string, error) {
	s.mu.RLock()
	ch, ok := s.channels[name]
	s.mu.RUnlock()

	if !ok {
		return nil, errors.New("channel not found")
	}

	select {
	case val := <-ch:
		return val, nil
	default:
		return nil, errors.New("channel is empty")
	}
}
