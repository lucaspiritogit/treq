package modal

import "github.com/rivo/tview"

func ShowErrorModal(app *tview.Application, appFlex *tview.Flex, errorMsg string) {
	modal := tview.NewModal().
		SetText(errorMsg).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			app.SetRoot(appFlex, true)
			app.SetFocus(appFlex)
		})

	pages := tview.NewPages().
		AddPage("background", appFlex, true, true).
		AddPage("error", modal, true, true)

	app.SetRoot(pages, true)
	app.SetFocus(modal)
}
