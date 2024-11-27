package response

import (
	"encoding/json"
	"treq/internal/ui/components/basecomponent"
	"treq/internal/ui/state"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ResponseTextView struct {
	basecomponent.BaseComponent
	View *tview.TextView
}

func NewResponseTextView(app *tview.Application, state *state.AppState) *ResponseTextView {
	return &ResponseTextView{
		BaseComponent: *basecomponent.NewBaseComponent(app, state),
		View:          tview.NewTextView(),
	}
}

func (r *ResponseTextView) Build(){
	r.View.SetDynamicColors(true).
		SetRegions(true).
		SetWrap(true).
		SetScrollable(true).
		SetBorder(true).
		SetTitle("Response")

	r.View.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyPgUp:
			r.View.ScrollToBeginning()
			return nil
		case tcell.KeyPgDn:
			r.View.ScrollToEnd()
			return nil
		case tcell.KeyEsc:
			r.State.FocusAppFlexContainer()
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
