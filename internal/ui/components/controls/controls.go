package controls

import "github.com/rivo/tview"

func GetControlsTextView() *tview.TextView {
	controlsTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetText(`
[blue]ctrl+k[white] Toggle Keyboard Shortcuts

[blue]k[white] Focus Keyboard Shortcuts

[red]Esc[white] Exit Input mode

HTTP Verbs:
[red]g[white]   GET
[red]p[white]   POST
[red]t[white]   PUT
[red]d[white]   DELETE

Navigation:
[red]i[white]   Focus URL
[red]r[white]   Focus response
[red]q[white]   Quit`).
		SetRegions(true)

	return controlsTextView
}
