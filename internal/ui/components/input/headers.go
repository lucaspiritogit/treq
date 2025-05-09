package input

import (
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
}

func (h *Headers) SetAppFlex(appFlex *tview.Flex) {
	h.appFlex = appFlex
}

func NewHeaders() *Headers {
	headerContainer := tview.NewFlex().SetDirection(tview.FlexRow)
	headerContainer.SetTitle("Headers").SetTitleAlign(tview.AlignLeft)
	headerContainer.SetBorder(true)

	rowContainer := tview.NewFlex().SetDirection(tview.FlexRow)
	rowContainer.SetBorderPadding(1, 1, 0, 0)
	headerContainer.AddItem(rowContainer, 0, 1, false)

	headers := &Headers{
		Container:   headerContainer,
		Rows:        rowContainer,
		Fields:      []*tview.InputField{},
		HeaderPages: make(map[int][]map[string]string),
	}

	return headers
}

func (h *Headers) GetHeadersContainer(app *tview.Application, appState *state.AppState) *tview.Flex {
	h.addNewHeader(app, appState)

	h.Container.SetFocusFunc(func() {
		appState.SetInputActive(true)
		app.SetFocus(h.Rows.GetItem(0))
	})

	h.Container.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEsc:
			h.saveCurrentPageHeaders()
			app.SetFocus(h.appFlex)
			appState.SetInputActive(false)
			return nil
		}
		return event
	})

	return h.Container
}


func (h *Headers) GetHeaders() []map[string]string {
	h.saveCurrentPageHeaders()

	allHeaders := []map[string]string{}
	return allHeaders
}

func (h *Headers) addNewHeader(app *tview.Application, appState *state.AppState) {
	h.AddHeadersRow(app, "", "", appState)
}

func (h *Headers) AddHeadersRow(app *tview.Application, key, value string, appState *state.AppState) {
	inputKey := tview.NewInputField().SetFieldWidth(30).SetLabel("Key:").SetText(key)
	inputKey.SetFieldBackgroundColor(tcell.ColorWhite)

	inputValue := tview.NewInputField().SetLabel("Value:").SetText(value)
	inputValue.SetFieldBackgroundColor(tcell.ColorWhite)

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
		if event.Key() == tcell.KeyTab || event.Key() == tcell.KeyRight {
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
	row.AddItem(inputKey, 0, 1, true)
	row.AddItem(inputValue, 0, 1, false)
	h.Rows.AddItem(row, 2, 1, false)
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
