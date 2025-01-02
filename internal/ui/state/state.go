package state

import (
	"sync"

	"github.com/rivo/tview"
)

type AppState struct {
	isInputActive    bool
	mu               sync.Mutex
	AppFlexContainer *tview.Flex
	App              *tview.Application
}

func NewAppState(appFlexContainer *tview.Flex, app *tview.Application) *AppState {
	return &AppState{
		isInputActive:    false,
		AppFlexContainer: appFlexContainer,
		App:              app,
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

func (s *AppState) FocusAppFlexContainer() {
	s.SetInputActive(false)
	s.App.SetFocus(s.AppFlexContainer)
}
