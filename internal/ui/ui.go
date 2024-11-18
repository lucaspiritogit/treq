package ui

import (
	"log"
	"time"
	"treq/internal/constants"
	"treq/internal/models"
	"treq/internal/repository"
	"treq/internal/service"
	"treq/internal/ui/components/controls"
	"treq/internal/ui/components/dropdown"
	inputfield "treq/internal/ui/components/input"
	"treq/internal/ui/components/request"
	"treq/internal/ui/components/response"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

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

func InitializeAppUI(app *tview.Application, requestList *tview.List, requestService *service.RequestService, requestRepository *repository.RequestRepository) *tview.Flex {
	isKeyboardControlTextViewRemove := false

	httpVerbDropdown := dropdown.GetHttpVerbDropdown()
	responseTextView := response.GetResponseTextView()
	responseMetadata := response.GetResponseMetadataTextView()
	urlInputField := inputfield.NewURLInputField(responseTextView, responseMetadata, httpVerbDropdown, app)
	bodyRequestTextView := request.GetRequestTextView()
	controlsTextView := controls.GetControlsTextView()

	controlsTextView.SetBorder(true).SetTitle("Keyboard Shortcuts")
	inputArea := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(httpVerbDropdown, 1, 0, false).
		AddItem(urlInputField, 1, 0, false)

	appFlex := tview.NewFlex().
		AddItem(requestList, 25, 0, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(inputArea, 3, 1, false).
			AddItem(responseMetadata, 1, 1, false).
			AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(bodyRequestTextView, 0, 1, false).
				AddItem(responseTextView, 0, 1, false), 0, 8, false),
			0, 1, false).
		AddItem(controlsTextView, 20, 0, false)

	urlInputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			app.SetFocus(appFlex)
			return nil
		}
		return event
	})

	responseTextView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			app.SetFocus(appFlex)
			return nil
		}
		return event
	})

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		saveRequestInput := tview.NewInputField().
			SetLabel("Request Name: ").
			SetFieldWidth(20)

		switch event.Key() {
		case tcell.KeyCtrlK:
			if isKeyboardControlTextViewRemove {
				appFlex.AddItem(controlsTextView, 20, 0, false)
			} else {
				appFlex.RemoveItem(controlsTextView)
			}
			isKeyboardControlTextViewRemove = !isKeyboardControlTextViewRemove
			return nil
		case tcell.KeyCtrlS:
			app.SetRoot(saveRequestInput, true).SetFocus(saveRequestInput)
			app.SetInputCapture(nil)

			saveRequestInput.SetDoneFunc(func(key tcell.Key) {
				if key == tcell.KeyEnter {
					requestName := saveRequestInput.GetText()
					_, httpMethodCurrentOption := httpVerbDropdown.GetCurrentOption()

					savedRequest := models.SavedRequest{
						Method:    httpMethodCurrentOption,
						Title:     requestName,
						URL:       urlInputField.GetText(),
						Body:      responseTextView.GetText(false),
						CreatedAt: time.Now(),
					}
					err := requestRepository.SaveRequest(savedRequest)
					if err != nil {
						log.Println("Error saving request:", err)
						return
					}
					requestService.RefreshRequestsList(requestList, requestRepository)
					app.SetRoot(appFlex, true)
					app.SetFocus(requestList)
				} else if key == tcell.KeyEsc {
					app.SetRoot(appFlex, true)
					app.SetFocus(requestList)
				}
			})
			return nil
		case tcell.KeyLeft:
			if app.GetFocus() == responseTextView || app.GetFocus() == bodyRequestTextView || app.GetFocus() == urlInputField {
				return nil
			}
			app.SetFocus(requestList)
		}
		switch event.Rune() {
		case 'i':
			app.SetFocus(urlInputField)
			return nil
		case 'r':
			app.SetFocus(responseTextView)
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
			app.SetFocus(controlsTextView)
			return nil
		case 'q':
			app.Stop()
			return nil
		}
		return event
	})

	requestList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			app.SetFocus(appFlex)
			return nil
		}
		if event.Key() == tcell.KeyEnter {
			index := requestList.GetCurrentItem()
			savedRequest := requestRepository.GetRequestById(index)
			methodToIndex := map[string]int{
				"GET":    constants.HTTPVerbGetDropdownIndex,
				"POST":   constants.HTTPVerbPostDropwdownIndex,
				"PUT":    constants.HTTPVerbPutDropdownIndex,
				"DELETE": constants.HTTPVerbDeleteDropdownIndex,
			}

			urlInputField.SetText(savedRequest.URL)
			bodyRequestTextView.SetText("body: " + savedRequest.Body)
			if index, ok := methodToIndex[savedRequest.Method]; ok {
				httpVerbDropdown.SetCurrentOption(index)
			}
			app.SetFocus(urlInputField)

			return nil
		}
		if event.Key() == tcell.KeyBackspace || event.Key() == tcell.KeyDelete {
			index := requestList.GetCurrentItem()
			go func() {
				requestRepository.DeleteRequestById(index)

				app.QueueUpdateDraw(func() {
					requestService.RefreshRequestsList(requestList, requestRepository)
				})
			}()
			return nil
		}
		return event
	})

	return appFlex
}
