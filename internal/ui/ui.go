package ui

import (
	"treq/internal/constants"
	"treq/internal/repository"
	"treq/internal/service"
	"treq/internal/ui/components/controls"
	"treq/internal/ui/components/dropdown"
	"treq/internal/ui/components/input"
	"treq/internal/ui/components/modal"
	"treq/internal/ui/components/request"
	"treq/internal/ui/components/response"
	"treq/internal/ui/components/shortcuts"
	"treq/internal/ui/state"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func InitializeAppUI(app *tview.Application, requestList *tview.TreeView, requestService *service.RequestService, requestRepository *repository.RequestRepository) *tview.Flex {

	spacer := tview.NewBox()
	appFlexContainer := tview.NewFlex().SetDirection(tview.FlexRow)
	appState := state.NewAppState(appFlexContainer, app)

	httpVerbDropdown := dropdown.GetHttpVerbDropdown()
	responseTextView := response.NewResponseTextView(appState)
	responseMetadata := response.NewResponseMetadataTextView()
	requestBody := request.NewRequestBody(app, appState, appFlexContainer)
	controlsTextView := controls.NewControlsModal()
	headers := input.NewHeaders()
	headersContainer := headers.GetHeadersContainer(app, appState)

	urlInputField := input.NewURLInputField(responseTextView.TextView, responseMetadata.View, httpVerbDropdown, headers, app)
	inputArea := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(httpVerbDropdown, 7, 0, false).
		AddItem(urlInputField, 0, 1, false)

	shortcutsTextView := shortcuts.NewShortcutsView(app)

	appFlexContainer.
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(requestList, 25, 0, false).
				AddItem(
					tview.NewFlex().SetDirection(tview.FlexRow).
						AddItem(inputArea, 3, 0, false).
						AddItem(headersContainer, 18, 1, false).
						AddItem(spacer, 0, 1, false).
						AddItem(responseMetadata.View, 1, 0, false).
						AddItem(
							tview.NewFlex().SetDirection(tview.FlexColumn).
								AddItem(requestBody.TextArea, 0, 1, false).
								AddItem(responseTextView.TextView, 0, 1, false),
							0, 8, false),
					0, 1, false),
			0, 1, false).
		AddItem(shortcutsTextView.View, 2, 0, false)

	headers.SetAppFlex(appFlexContainer)


	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if appState.IsInputActive() {
			return event
		}

		requestNameInputField := tview.NewInputField().
			SetLabel("Request Name: ").
			SetFieldWidth(20)

		switch event.Key() {
		case tcell.KeyEsc:
			appState.SetInputActive(false)
			app.SetRoot(appFlexContainer, true)
			app.SetFocus(appFlexContainer)
			return nil
		case tcell.KeyCtrlN:
			app.SetRoot(requestNameInputField, true).SetFocus(requestNameInputField)
			appState.SetInputActive(true)

			requestNameInputField.SetDoneFunc(func(key tcell.Key) {
				if key == tcell.KeyEnter {
					appState.SetInputActive(false)
					requestName := requestNameInputField.GetText()

					if requestName == "" {
						modal.ShowErrorModal(app, appFlexContainer, "Provide a name to save the request.")
						return
					}
					app.SetRoot(appFlexContainer, true)
					app.SetFocus(requestList)
				} else if key == tcell.KeyEsc {
					app.SetRoot(appFlexContainer, true)
					app.SetFocus(requestList)
				}
			})
		case tcell.KeyLeft:
			if app.GetFocus() == responseTextView.TextView || app.GetFocus() == requestBody.TextArea || app.GetFocus() == urlInputField {
				return nil
			}
			app.SetFocus(requestList)
		}
		switch event.Rune() {
		case 'i':
			appState.SetInputActive(true)
			app.SetFocus(urlInputField)
			return nil
		case 'r':
			app.SetFocus(responseTextView.TextView)
			return nil
		case 'p':
			httpVerbDropdown.SetCurrentOption(constants.HTTPVerbPostDropwdownIndex)
			return nil
		case 'g':
			httpVerbDropdown.SetCurrentOption(constants.HTTPVerbGetDropdownIndex)
		case 'd':
			httpVerbDropdown.SetCurrentOption(constants.HTTPVerbDeleteDropdownIndex)
			return nil
		case 'e':
			httpVerbDropdown.SetCurrentOption(constants.HTTPVerbPutDropdownIndex)
			return nil
		case 'k':
			app.SetFocus(controlsTextView.View)
			return nil
		case 'b':
			appState.SetInputActive(true)
			app.SetFocus(requestBody.TextArea)
			return nil
		case 'h':
			app.SetFocus(headers.Container)
			appState.SetInputActive(true)
			return nil
		case 'q':
			app.Stop()
			return nil
		case '?':
			app.SetRoot(controlsTextView.View, true).SetFocus(controlsTextView.View)
			return nil
		}
		return event
	})

	return appFlexContainer
}