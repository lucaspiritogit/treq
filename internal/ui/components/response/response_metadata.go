package response

import "github.com/rivo/tview"

func GetResponseMetadataTextView() *tview.TextView {
	responseMetadataTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWrap(false).
		SetScrollable(true)

	responseMetadataTextView.SetText("[white]Status: [green]- [white]Content-Length: [green]0")
	responseMetadataTextView.SetTextAlign(tview.AlignRight)

	return responseMetadataTextView
}
