package request

import "github.com/rivo/tview"

func GetRequestTextView() *tview.TextView {
	requestTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetScrollable(true)
	requestTextView.SetBorder(true).SetTitle("Request")

	return requestTextView
}
