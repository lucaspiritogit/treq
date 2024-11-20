package state

import "sync"

type AppState struct {
	isInputActive bool
	mu            sync.Mutex
}

func NewAppState() *AppState {
	return &AppState{
		isInputActive: false,
	}
}

func (s *AppState) SetInputActive(active bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.isInputActive = active
}

func (s *AppState) IsInputActive() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.isInputActive
}
