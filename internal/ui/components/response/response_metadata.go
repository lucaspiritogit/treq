package response

import "github.com/rivo/tview"

type ResponseMetadataTextView struct {
	View *tview.TextView
}

func NewResponseMetadataTextView() *ResponseMetadataTextView {
	view := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetScrollable(true)

	view.SetText("[white]Status: [green]- [white]Content-Length: [green]0")
	view.SetTextAlign(tview.AlignRight)

	r := &ResponseMetadataTextView{View: view}
	return r
}