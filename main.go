package main

import (
	"log"
	"treq/internal/storage"
	"treq/internal/ui/components/controls"
	httpverbdropdown "treq/internal/ui/components/dropdown"
	urlinputfield "treq/internal/ui/components/input"
	"treq/internal/ui/components/request"
	"treq/internal/ui/components/response"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()
	isKeyboardControlTextViewRemove := false

	requestList := tview.NewList()

	requestManager, err := storage.NewRequestManager(requestList, app)
	if err != nil {
		panic(err)
	}
	storage.RefreshRequestsList(requestList, requestManager)

	responseMetadata := response.GetResponseMetadataTextView()

	httpVerbDropdown := httpverbdropdown.GetHttpVerbDropdown()
	responseTextView := response.GetResponseTextView()

	urlInputField := urlinputfield.NewURLInputField(responseTextView, responseMetadata, httpVerbDropdown, app)
	bodyRequestTextView := request.GetRequestTextView()

	controlsTextView := controls.GetControlsTextView()

	controlsTextView.SetBorder(true).SetTitle("Keyboard Shortcuts")

	responseTextView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyPgUp:
			responseTextView.ScrollToBeginning()
		case tcell.KeyPgDn:
			responseTextView.ScrollToEnd()
		}
		return event
	})

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

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			urlInputField.SetBorderColor(tcell.ColorWhite)
			app.SetFocus(responseTextView)
			return nil
		}

		switch event.Key() {
		case tcell.KeyCtrlK:
			if isKeyboardControlTextViewRemove {
				appFlex.AddItem(controlsTextView, 25, 0, false)
			} else {
				appFlex.RemoveItem(controlsTextView)
			}
			isKeyboardControlTextViewRemove = !isKeyboardControlTextViewRemove
			return nil
		case tcell.KeyCtrlS:
			go func() {
				_, httpMethodCurrentOption := httpVerbDropdown.GetCurrentOption()

				err := requestManager.SaveRequest(
					storage.SavedRequest{
						Method: httpMethodCurrentOption,
						URL:    urlInputField.GetText(),
						Body:   responseTextView.GetText(false),
					},
				)
				if err != nil {
					log.Println("Error saving request:", err)
					return
				}

				app.QueueUpdateDraw(func() {
					storage.RefreshRequestsList(requestList, requestManager)
				})
			}()
			return nil
		case tcell.KeyLeft:
			if app.GetFocus() == responseTextView || app.GetFocus() == bodyRequestTextView || app.GetFocus() == urlInputField {
				return nil
			}
			app.SetFocus(requestList)
		}

		switch event.Rune() {
		case 'i':
			urlInputField.SetBorderColor(tcell.ColorGreen)
			app.SetFocus(urlInputField)
			return nil
		case 'r':
			app.SetFocus(responseTextView)
			return nil
		case 'p':
			httpVerbDropdown.SetCurrentOption(1)
			return nil
		case 'g':
			httpVerbDropdown.SetCurrentOption(0)
		case 'd':
			httpVerbDropdown.SetCurrentOption(3)
			return nil
		case 'e':
			httpVerbDropdown.SetCurrentOption(2)
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
		if event.Key() == tcell.KeyEnter {
			index := requestList.GetCurrentItem()
			if index == 0 {
				return nil
			}
			title, secondary := requestList.GetItemText(index)
			urlInputField.SetText(title)
			bodyRequestTextView.SetText(secondary)
			return nil
		}
		if event.Key() == tcell.KeyCtrlD {
			index := requestList.GetCurrentItem()
			go func() {
				requestManager.DeleteRequestById(index)

				app.QueueUpdateDraw(func() {
					storage.RefreshRequestsList(requestList, requestManager)
				})
			}()
			return nil
		}
		return event
	})

	if err := app.SetRoot(appFlex, true).Run(); err != nil {
		panic(err)
	}
}
