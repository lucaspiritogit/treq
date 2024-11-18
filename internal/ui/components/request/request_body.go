package request

import "github.com/rivo/tview"

func GetRequestBody() *tview.TextView {
	requestTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetScrollable(true)
	requestTextView.SetBorder(true).SetTitle("Body")

	return requestTextView
}
