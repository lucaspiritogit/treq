package response

import (
	"encoding/json"
	"treq/internal/ui/state"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ResponseTextView struct {
	TextView       *tview.TextView
	AppState 			*state.AppState
}

func NewResponseTextView(appState *state.AppState) *ResponseTextView {
	view := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(true).
		SetScrollable(true)

	view.SetBorder(true).SetTitle("Response")

	r := &ResponseTextView{TextView: view, AppState: appState}
	r.setInputCapture()

	return r
}

func (r *ResponseTextView) setInputCapture() {
	r.TextView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyPgUp:
			r.TextView.ScrollToBeginning()
			return nil
		case tcell.KeyPgDn:
			r.TextView.ScrollToEnd()
			return nil
		case tcell.KeyEsc:
			r.AppState.FocusAppFlexContainer()
			return nil
		}
		return event
	})
}

func FormatJSON(input string) (string, error) {
	var prettyJSON any
	err := json.Unmarshal([]byte(input), &prettyJSON)
	if err != nil {
		return "", err
	}

	prettyJSONBytes, err := json.MarshalIndent(prettyJSON, "", "  ")
	if err != nil {
		return "", err
	}

	return string(prettyJSONBytes), nil
}
