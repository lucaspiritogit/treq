package request

import "github.com/rivo/tview"

func GetRequestBody() *tview.TextArea {
	requestBody := tview.NewTextArea().
		SetWrap(true)

	requestBody.SetTitle("Request Body").SetBorder(true)
	return requestBody
}
