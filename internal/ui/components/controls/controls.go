package controls

import (
	"github.com/rivo/tview"
)

type ControlsModal struct {
	View *tview.TextView
}

func NewControlsModal() *ControlsModal {
	controlsText := `
    [blue]?[white] Toggle Keyboard Shortcuts

    [blue]CTRL + S[white] Save Request

    [blue]k[white] Focus Keyboard Shortcuts

    [blue]Esc[white] Exit Input mode

    HTTP Verbs:
    [blue]g[white]   GET
    [blue]p[white]   POST
    [blue]e[white]   PUT
    [blue]d[white]   DELETE

    Navigation:
    [blue]i[white]   Focus URL
    [blue]r[white]   Focus response
    [blue]q[white]   Quit
    `

	controlsTextView := tview.NewTextView().
		SetText(controlsText).
		SetTextAlign(tview.AlignLeft).
		SetScrollable(true).
		SetDynamicColors(true)

	controlsTextView.SetBorder(true).SetTitle("Controls")

	return &ControlsModal{View: controlsTextView}
}
