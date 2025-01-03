package ui

import (
	"time"
	"treq/internal/constants"
	"treq/internal/models"
	"treq/internal/repository"
	"treq/internal/service"
	"treq/internal/ui/components/controls"
	"treq/internal/ui/components/dropdown"
	"treq/internal/ui/components/input"
	"treq/internal/ui/components/modal"
	"treq/internal/ui/components/request"
	"treq/internal/ui/components/response"
	"treq/internal/ui/state"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func InitializeAppUI(app *tview.Application, requestList *tview.List, requestService *service.RequestService, requestRepository *repository.RequestRepository) *tview.Flex {

	isControlsModalOpen := false
	spacer := tview.NewBox()
	appFlexContainer := tview.NewFlex().SetDirection(tview.FlexRow)
	appState := state.NewAppState(appFlexContainer, app)

	httpVerbDropdown := dropdown.GetHttpVerbDropdown()

	responseTextView := response.NewResponseTextView(app, appState)
	responseTextView.Build()

	responseMetadata := response.GetResponseMetadataTextView()
	requestBody := request.GetRequestBody()
	controlsTextView := controls.GetControlsModal()

	headers := input.NewHeaders()
	headersContainer := headers.GetHeadersContainer(app, appState)

	urlInputField := input.NewURLInputField(responseTextView.View, responseMetadata, httpVerbDropdown, headers, app)
	inputArea := tview.NewFlex().
		SetDirection(tview.FlexColumn).
		AddItem(httpVerbDropdown, 6, 0, false).
		AddItem(urlInputField, 0, 1, false)

	shortcutsTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetText("[blue]?[white] Show controls | [blue]Ctrl + s[white] Save req | [blue]q[white] Quit | [blue]i[white] URL | [blue]b[white] Body | [blue]r[white] Response body").
		SetScrollable(false).
		SetWrap(true)

	appFlexContainer.
		AddItem(
			tview.NewFlex().SetDirection(tview.FlexColumn).
				AddItem(requestList, 25, 0, false).
				AddItem(
					tview.NewFlex().SetDirection(tview.FlexRow).
						AddItem(inputArea, 1, 0, false).
						AddItem(headersContainer, 6, 1, false).
						AddItem(spacer, 0, 1, false).
						AddItem(responseMetadata, 1, 0, false).
						AddItem(
							tview.NewFlex().SetDirection(tview.FlexColumn).
								AddItem(requestBody, 0, 1, false).
								AddItem(responseTextView.View, 0, 1, false),
							0, 8, false),
					0, 1, false),
			0, 1, false).
		AddItem(shortcutsTextView, 2, 0, false)

	headers.SetAppFlex(appFlexContainer)

	urlInputField.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			appState.SetInputActive(false)
			app.SetFocus(appFlexContainer)
			return nil
		}
		return event
	})

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if appState.IsInputActive() {
			return event
		}

		requestNameInputField := tview.NewInputField().
			SetLabel("Request Name: ").
			SetFieldWidth(20)

		switch event.Key() {
		case tcell.KeyCtrlS:
			index := requestList.GetCurrentItem()
			existingRequest := requestRepository.GetRequestById(index)

			app.SetRoot(requestNameInputField, true).SetFocus(requestNameInputField)
			appState.SetInputActive(true)

			requestNameInputField.SetDoneFunc(func(key tcell.Key) {
				if key == tcell.KeyEnter {
					appState.SetInputActive(false)
					requestName := requestNameInputField.GetText()
					_, httpMethodCurrentOption := httpVerbDropdown.GetCurrentOption()

					if requestName == "" {
						modal.ShowErrorModal(app, appFlexContainer, "Provide a name to save the request.")
						return
					}

					savedRequest := models.SavedRequest{
						Method:    httpMethodCurrentOption,
						Title:     requestName,
						URL:       urlInputField.GetText(),
						Body:      requestBody.GetText(),
						CreatedAt: time.Now(),
					}
					if existingRequest.ID != 0 {
						requestRepository.UpdateRequest(savedRequest, index)
						requestService.RefreshRequestsList(requestList, requestRepository)
						savedHeaders := []models.SavedHeaders{}
						for page := 0; page < headers.TotalPages; page++ {
							for _, headerData := range headers.HeaderPages[page] {
								savedHeader := models.SavedHeaders{
									RequestId: existingRequest.ID,
									Key:       headerData["key"],
									Value:     headerData["value"],
									Page:      page,
									CreatedAt: time.Now(),
								}
								savedHeaders = append(savedHeaders, savedHeader)
							}
						}

						a := requestRepository.SaveHeaders(existingRequest.ID, savedHeaders)
						if a != nil {
							modal.ShowErrorModal(app, appFlexContainer, a.Error())
							return
						}

						app.SetRoot(appFlexContainer, true)
						app.SetFocus(requestList)
						return
					}
					createRequest(requestRepository, savedRequest, app, appFlexContainer, headers)
					requestService.RefreshRequestsList(requestList, requestRepository)
					app.SetRoot(appFlexContainer, true)
					app.SetFocus(requestList)
				} else if key == tcell.KeyEsc {
					app.SetRoot(appFlexContainer, true)
					app.SetFocus(requestList)
				}
			})
			return nil
		case tcell.KeyCtrlN:
			app.SetRoot(requestNameInputField, true).SetFocus(requestNameInputField)
			appState.SetInputActive(true)

			requestNameInputField.SetDoneFunc(func(key tcell.Key) {
				if key == tcell.KeyEnter {
					appState.SetInputActive(false)
					requestName := requestNameInputField.GetText()
					_, httpMethodCurrentOption := httpVerbDropdown.GetCurrentOption()

					if requestName == "" {
						modal.ShowErrorModal(app, appFlexContainer, "Provide a name to save the request.")
						return
					}

					savedRequest := models.SavedRequest{
						Method:    httpMethodCurrentOption,
						Title:     requestName,
						URL:       urlInputField.GetText(),
						Body:      requestBody.GetText(),
						CreatedAt: time.Now(),
					}
					createRequest(requestRepository, savedRequest, app, appFlexContainer, headers)
					requestService.RefreshRequestsList(requestList, requestRepository)
					app.SetRoot(appFlexContainer, true)
					app.SetFocus(requestList)
				} else if key == tcell.KeyEsc {
					app.SetRoot(appFlexContainer, true)
					app.SetFocus(requestList)
				}
			})
		case tcell.KeyLeft:
			if app.GetFocus() == responseTextView.View || app.GetFocus() == requestBody || app.GetFocus() == urlInputField {
				return nil
			}
			app.SetFocus(requestList)
		}
		switch event.Rune() {
		case '?':
			if isControlsModalOpen {
				app.SetRoot(appFlexContainer, true)
			} else {
				app.SetRoot(controlsTextView, true)
			}
			isControlsModalOpen = !isControlsModalOpen
			return nil
		case 'i':
			appState.SetInputActive(true)
			app.SetFocus(urlInputField)
			return nil
		case 'r':
			app.SetFocus(responseTextView.View)
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
		case 'b':
			appState.SetInputActive(true)
			app.SetFocus(requestBody)
			return nil
		case 'h':
			app.SetFocus(headers.Container)
			return nil
		case 'q':
			app.Stop()
			return nil
		}
		return event
	})

	requestBody.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			appState.SetInputActive(false)
			app.SetFocus(appFlexContainer)
			return nil
		}
		return event
	})

	requestList.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			app.SetFocus(appFlexContainer)
			return nil
		}
		if event.Key() == tcell.KeyEnter {
			headers.HeaderPages = map[int][]map[string]string{}
			index := requestList.GetCurrentItem()
			existingRequest := requestRepository.GetRequestById(index)
			methodToIndex := map[string]int{
				"GET":    constants.HTTPVerbGetDropdownIndex,
				"POST":   constants.HTTPVerbPostDropwdownIndex,
				"PUT":    constants.HTTPVerbPutDropdownIndex,
				"DELETE": constants.HTTPVerbDeleteDropdownIndex,
			}

			savedHeaders := requestRepository.GetHeadersByRequestId(existingRequest.ID)
			for page, pageHeaders := range savedHeaders {
				if len(pageHeaders) > 0 {
					if headers.HeaderPages[page] == nil {
						headers.HeaderPages[page] = []map[string]string{}
						headers.AddHeadersRow(app, "", "", appState)
					}

					for _, header := range pageHeaders {
						headers.HeaderPages[page] = append(headers.HeaderPages[page], map[string]string{
							"key":   header.Key,
							"value": header.Value,
						})
					}
				}
			}

			headers.TotalPages = len(savedHeaders)
			headers.UpdatePageDisplay(app, appState)
			urlInputField.SetText(existingRequest.URL)
			requestBody.SetText(existingRequest.Body, false)

			if index, ok := methodToIndex[existingRequest.Method]; ok {
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

	return appFlexContainer
}

func createRequest(requestRepository *repository.RequestRepository, savedRequest models.SavedRequest, app *tview.Application, appFlexContainer *tview.Flex, headers *input.Headers) {
	result, err := requestRepository.SaveRequest(savedRequest)
	if err != nil {
		modal.ShowErrorModal(app, appFlexContainer, err.Error())
		return
	}
	insertedId, err := result.LastInsertId()
	if err != nil {
		modal.ShowErrorModal(app, appFlexContainer, err.Error())
		return
	}
	requestId := int(insertedId)

	savedHeaders := []models.SavedHeaders{}
	for page := 0; page < headers.TotalPages; page++ {
		for _, headerData := range headers.HeaderPages[page] {
			savedHeader := models.SavedHeaders{
				RequestId: requestId,
				Key:       headerData["key"],
				Value:     headerData["value"],
				Page:      page,
				CreatedAt: time.Now(),
			}
			savedHeaders = append(savedHeaders, savedHeader)
		}
	}

	a := requestRepository.SaveHeaders(requestId, savedHeaders)
	if a != nil {
		modal.ShowErrorModal(app, appFlexContainer, a.Error())
		return
	}
}
