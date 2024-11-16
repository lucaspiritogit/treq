package httpverbdropdown

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func GetHttpVerbDropdown() *tview.DropDown {
	httpVerbDropdown := tview.NewDropDown().
		SetLabel("HTTP Verb: ").
		SetOptions([]string{"GET", "POST", "PUT", "DELETE"}, nil).
		SetCurrentOption(0)

	httpVerbDropdown.SetFieldBackgroundColor(tcell.ColorBlack)
	httpVerbDropdown.SetFieldTextColor(tcell.ColorGreen)
	httpVerbDropdown.SetSelectedFunc(func(text string, index int) {
		switch text {
		case "GET":
			httpVerbDropdown.SetFieldTextColor(tcell.ColorGreen)
		case "POST":
			httpVerbDropdown.SetFieldTextColor(tcell.ColorYellow)
		case "PUT":
			httpVerbDropdown.SetFieldTextColor(tcell.ColorBlue)
		case "DELETE":
			httpVerbDropdown.SetFieldTextColor(tcell.ColorRed)
		}
	})

	return httpVerbDropdown
}
