package input

import "github.com/rivo/tview"

func GetHeadersInputField() *tview.InputField {
	return tview.NewInputField().
		SetLabel("Headers: ").
		SetFieldWidth(55).
		SetAcceptanceFunc(tview.InputFieldMaxLength(1024)).
		SetText("")
}
