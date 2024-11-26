package input

import (
	"fmt"
	"strings"

	"treq/internal/ui/state"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Headers struct {
	Container   *tview.Flex
	Rows        *tview.Flex
	Fields      []*tview.InputField
	appFlex     *tview.Flex
	HeaderPages map[int][]map[string]string
	currentPage int
	pageSize    int
	TotalPages  int
	pageLabel   *tview.TextView
}

func (h *Headers) SetAppFlex(appFlex *tview.Flex) {
	h.appFlex = appFlex
}

func NewHeaders() *Headers {
	headerContainer := tview.NewFlex().SetDirection(tview.FlexRow)
	headerContainer.SetTitle("Headers").SetTitleAlign(tview.AlignLeft)
	headerContainer.SetBorder(true)

	rowContainer := tview.NewFlex().SetDirection(tview.FlexRow)

	pageLabel := tview.NewTextView().
		SetTextColor(tcell.ColorWhite).
		SetText("Page 1/1")

	headerContainer.AddItem(rowContainer, 0, 1, false)
	headerContainer.AddItem(pageLabel, 1, 0, false)

	headers := &Headers{
		Container:   headerContainer,
		Rows:        rowContainer,
		Fields:      []*tview.InputField{},
		HeaderPages: make(map[int][]map[string]string),
		currentPage: 0,
		pageSize:    3,
		TotalPages:  1,
		pageLabel:   pageLabel,
	}

	headers.addNewHeader(nil, nil)

	return headers
}

func (h *Headers) GetHeadersContainer(app *tview.Application, appState *state.AppState) *tview.Flex {
	h.addNewHeader(app, appState)

	h.Container.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			appState.SetInputActive(true)
			app.SetFocus(h.Rows.GetItem(0))
			return nil

		case tcell.KeyEsc:
			h.saveCurrentPageHeaders()
			app.SetFocus(h.appFlex)
			appState.SetInputActive(false)
			return nil

		case tcell.KeyLeft:
			h.prevPage(app, appState)
			return nil

		case tcell.KeyRight:
			h.nextPage(app, appState)
			return nil
		}
		return event
	})

	return h.Container
}

func (h *Headers) prevPage(app *tview.Application, appState *state.AppState) {
	h.saveCurrentPageHeaders()
	if h.currentPage > 0 {
		h.currentPage--
		h.UpdatePageDisplay(app, appState)
	}
}

func (h *Headers) nextPage(app *tview.Application, appState *state.AppState) {
	h.saveCurrentPageHeaders()
	if h.currentPage < h.TotalPages-1 {
		h.currentPage++
		h.UpdatePageDisplay(app, appState)
	}
}

func (h *Headers) UpdatePageDisplay(app *tview.Application, appState *state.AppState) {
	h.Rows.Clear()

	currentPageHeaders := h.HeaderPages[h.currentPage]

	for _, header := range currentPageHeaders {
		h.AddHeadersRow(app, header["key"], header["value"], appState)
	}

	h.pageLabel.SetText(fmt.Sprintf("Page %d/%d", h.currentPage+1, h.TotalPages))

	if h.Rows.GetItemCount() > 0 {
		firstRow := h.Rows.GetItem(0).(*tview.Flex)
		app.SetFocus(firstRow.GetItem(0).(*tview.InputField))
	}
}

func (h *Headers) GetHeaders() []map[string]string {
	h.saveCurrentPageHeaders()

	allHeaders := []map[string]string{}
	for page := 0; page < h.TotalPages; page++ {
		allHeaders = append(allHeaders, h.HeaderPages[page]...)
	}
	return allHeaders
}

func (h *Headers) addNewHeader(app *tview.Application, appState *state.AppState) {
	if len(h.HeaderPages[h.currentPage]) >= h.pageSize {
		h.currentPage = h.TotalPages
		h.TotalPages++
		h.HeaderPages[h.currentPage] = []map[string]string{}
	}

	h.AddHeadersRow(app, "", "", appState)
}

func (h *Headers) AddHeadersRow(app *tview.Application, key, value string, appState *state.AppState) {
	inputKey := tview.NewInputField().SetFieldWidth(30).SetLabel("Key:").SetText(key)
	inputValue := tview.NewInputField().SetLabel("Value:").SetText(value)

	navigateUp := func(currentField *tview.InputField) {
		fields := h.getAllInputFields()
		for i, field := range fields {
			if field == currentField && i > 0 {
				app.SetFocus(fields[i-1])
				break
			}
		}
	}

	navigateDown := func(currentField *tview.InputField) {
		fields := h.getAllInputFields()
		for i, field := range fields {
			if field == currentField && i > 0 && i < len(fields)-1 {
				app.SetFocus(fields[i+1])
				break
			}
		}
	}

	inputKey.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab || event.Key() == tcell.KeyDown {
			app.SetFocus(inputValue)
			return nil
		}

		if event.Key() == tcell.KeyUp {
			navigateUp(inputKey)
			return nil
		}

		if event.Key() == tcell.KeyDown {
			navigateDown(inputKey)
			return nil
		}

		if event.Key() == tcell.KeyEsc {
			app.SetFocus(h.appFlex)
			return nil
		}

		return event
	})

	inputValue.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			currentRowIndex := h.getRowIndex(inputValue)
			h.saveCurrentPageHeaders()

			if currentRowIndex == h.Rows.GetItemCount()-1 {
				h.addNewHeader(app, appState)
			}

			if currentRowIndex+1 < h.Rows.GetItemCount() {
				nextRow := h.Rows.GetItem(currentRowIndex + 1).(*tview.Flex)
				nextKey := nextRow.GetItem(0).(*tview.InputField)
				app.SetFocus(nextKey)
			}
			return nil
		}

		if event.Key() == tcell.KeyUp {
			navigateUp(inputValue)
			return nil
		}

		if event.Key() == tcell.KeyDown {
			navigateDown(inputValue)
			return nil
		}

		if event.Key() == tcell.KeyEsc {
			app.SetFocus(h.appFlex)
			return nil
		}

		return event
	})

	row := tview.NewFlex().SetDirection(tview.FlexColumn)
	row.AddItem(inputKey, 0, 2, true)
	row.AddItem(inputValue, 0, 3, false)
	h.Rows.AddItem(row, 0, 1, false)

	// h.HeaderPages[h.currentPage] = append(h.HeaderPages[h.currentPage],
	// 	map[string]string{"key": inputKey.GetText(), "value": inputValue.GetText()})

	h.pageLabel.SetText(fmt.Sprintf("Page %d/%d", h.currentPage+1, h.TotalPages))
}

func (h *Headers) saveCurrentPageHeaders() {
	updatedHeaders := []map[string]string{}
	for i := 0; i < h.Rows.GetItemCount(); i++ {
		row := h.Rows.GetItem(i).(*tview.Flex)
		inputKey := row.GetItem(0).(*tview.InputField)
		inputValue := row.GetItem(1).(*tview.InputField)

		key := strings.TrimSpace(inputKey.GetText())
		value := strings.TrimSpace(inputValue.GetText())

		if key != "" || value != "" {
			updatedHeaders = append(updatedHeaders,
				map[string]string{"key": key, "value": value})
		}
	}

	h.HeaderPages[h.currentPage] = updatedHeaders
}

func (h *Headers) getRowIndex(field *tview.InputField) int {
	for i := 0; i < h.Rows.GetItemCount(); i++ {
		row := h.Rows.GetItem(i).(*tview.Flex)
		if row.GetItem(0) == field || row.GetItem(1) == field {
			return i
		}
	}
	return -1
}

func (h *Headers) getAllInputFields() []*tview.InputField {
	fields := []*tview.InputField{}
	for i := 0; i < h.Rows.GetItemCount(); i++ {
		row := h.Rows.GetItem(i).(*tview.Flex)
		keyField := row.GetItem(0).(*tview.InputField)
		valueField := row.GetItem(1).(*tview.InputField)
		fields = append(fields, keyField, valueField)
	}
	return fields
}
