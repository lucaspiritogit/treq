package basecomponent

import (
	"treq/internal/ui/state"

	"github.com/rivo/tview"
)

type BaseComponent struct {
	App              *tview.Application
	AppFlexContainer *tview.Flex
	State            *state.AppState
}

func NewBaseComponent(app *tview.Application, state *state.AppState) *BaseComponent {
	return &BaseComponent{App: app, State: state}
}

type Component interface {
	Build()
}
