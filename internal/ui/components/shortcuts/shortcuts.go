package shortcuts

import (
	"github.com/rivo/tview"
)

type ShortcutsView struct {
	View *tview.TextView
}

func NewShortcutsView(app *tview.Application) *ShortcutsView {
	view := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[blue]?[white] Show controls | [blue]Ctrl + s[white] Save req | [blue]q[white] Quit | [blue]i[white] URL | [blue]b[white] Body | [blue]r[white] Response body").
		SetScrollable(false).
		SetWrap(true)
    shortcutView := &ShortcutsView{View: view}

	return shortcutView
}
