package request

import (
	"treq/internal/ui/state"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type RequestBody struct {
	TextArea *tview.TextArea
	appState *state.AppState
	app        *tview.Application
	appFlexContainer *tview.Flex
}

func NewRequestBody(app *tview.Application, appState *state.AppState, appFlexContainer *tview.Flex) *RequestBody {
	requestBody := &RequestBody{
		TextArea: tview.NewTextArea(),
		appState: appState,
		app:        app,
		appFlexContainer: appFlexContainer,
	}

	requestBody.TextArea.SetTitle("Body").SetBorder(true)
	requestBody.TextArea.SetWrap(true)
	requestBody.TextArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			requestBody.appState.SetInputActive(false)
			requestBody.app.SetFocus(requestBody.appFlexContainer)
			return nil
		}
		return event
	})

	return requestBody
}