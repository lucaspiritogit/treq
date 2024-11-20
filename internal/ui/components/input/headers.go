package input

import (
	"treq/internal/ui/state"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Headers struct {
	Container *tview.Flex
	Rows      *tview.Flex
	Fields    []*tview.InputField
	appFlex   *tview.Flex
}

func (h *Headers) SetAppFlex(appFlex *tview.Flex) {
	h.appFlex = appFlex
}

func NewHeaders() *Headers {
	headerContainer := tview.NewFlex().SetDirection(tview.FlexRow)
	headerContainer.SetTitle("Headers").SetTitleAlign(tview.AlignLeft)
	headerContainer.SetBorder(true)

	rowContainer := tview.NewFlex().SetDirection(tview.FlexRow)
	headerContainer.AddItem(rowContainer, 0, 1, false)

	return &Headers{
		Container: headerContainer,
		Rows:      rowContainer,
	}
}

func (h *Headers) GetHeadersContainer(app *tview.Application, appState *state.AppState) *tview.Flex {
	h.addNewHeader(app, appState)
	h.Container.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			appState.SetInputActive(true)
			app.SetFocus(h.Rows.GetItem(0))
			return nil
		}

		if event.Key() == tcell.KeyEsc {
			app.SetFocus(h.appFlex)
			appState.SetInputActive(false)
			return nil
		}
		return event
	})
	return h.Container
}

func (h *Headers) GetHeaders() []map[string]string {
	arr := []map[string]string{}
	for i := 0; i < h.Rows.GetItemCount(); i++ {
		row := h.Rows.GetItem(i).(*tview.Flex)

		inputKey := row.GetItem(0).(*tview.InputField)
		inputValue := row.GetItem(1).(*tview.InputField)
		key := inputKey.GetText()
		value := inputValue.GetText()

		arr = append(arr, map[string]string{"key": key, "value": value})
	}

	return arr
}

func (h *Headers) addNewHeader(app *tview.Application, appState *state.AppState) {
	inputKey := tview.NewInputField().SetFieldWidth(30).SetLabel("Key:")
	inputValue := tview.NewInputField().SetLabel("Value:")

	navigateUp := func(currentField *tview.InputField) {
		fields := h.getAllInputFields()
		for i, field := range fields {
			if field == currentField && i > 0 {
				app.SetFocus(fields[i-1])
				break
			}
		}
	}
	inputKey.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTab {
			app.SetFocus(inputValue)
			return nil
		}

		if event.Key() == tcell.KeyUp {
			navigateUp(inputKey)
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
