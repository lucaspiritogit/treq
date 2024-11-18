package inputfield

import (
	"fmt"
	"regexp"
	"treq/internal/http_request"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type URLInputField struct {
	*tview.InputField
	responseView     *tview.TextView
	responseMetadata *tview.TextView
	httpVerbDropdown *tview.DropDown
	app              *tview.Application
	isLoading        bool
}

func NewURLInputField(responseView, responseMetadata *tview.TextView, httpVerbDropdown *tview.DropDown, app *tview.Application) *URLInputField {
	urlField := &URLInputField{
		InputField:       tview.NewInputField(),
		responseView:     responseView,
		responseMetadata: responseMetadata,
		httpVerbDropdown: httpVerbDropdown,
		app:              app,
	}

	urlField.
		SetLabel("URL: ").
		SetFieldWidth(55).
		SetAcceptanceFunc(tview.InputFieldMaxLength(1024)).
		SetText("https://jsonplaceholder.typicode.com/todos")

	urlField.SetFocusFunc(func() {
	})

	urlField.SetDoneFunc(urlField.handleKeyPress)
	return urlField
}

func (u *URLInputField) handleKeyPress(key tcell.Key) {
	switch key {
	case tcell.KeyEnter:
		if !u.isLoading {
			u.executeRequest()
		}
		u.app.SetFocus(u.responseView)
	case tcell.KeyEsc:
		u.SetBorderColor(tcell.ColorWhite)
		u.app.SetFocus(u.responseView)
	}
}

func (u *URLInputField) executeRequest() {
	u.isLoading = true
	u.SetBorderColor(tcell.ColorYellow)
	u.responseMetadata.SetText("[yellow]Loading...")

	_, verb := u.httpVerbDropdown.GetCurrentOption()

	go func() {
		result := http_request.FetchUrl(u.GetText(), verb)

		u.app.QueueUpdateDraw(func() {
			formattedBody := colorizeJSON(result.Body)
			u.responseView.SetText(formattedBody)
			statusText := u.formatStatusCode(result.Resp.StatusCode)
			contentLength := len(result.Body)
			u.responseMetadata.SetText(fmt.Sprintf("%s [white]Content-Length: [blue]%d",
				statusText, contentLength))

			u.isLoading = false
			u.SetBorderColor(tcell.ColorWhite)

		})
	}()
}

func colorizeJSON(jsonStr string) string {
	keyRegex := regexp.MustCompile(`"(\w+)":`)
	colored := keyRegex.ReplaceAllString(jsonStr, `[yellow]"$1":[-]`)

	return colored
}

func (u *URLInputField) formatStatusCode(code int) string {
	switch {
	case code >= 200 && code < 300:
		return fmt.Sprintf("[white]Status: [green]%d", code)
	case code >= 300 && code < 400:
		return fmt.Sprintf("[white]Status: [yellow]%d", code)
	default:
		return fmt.Sprintf("[white]Status: [red]%d", code)
	}
}

func (u *URLInputField) Draw(screen tcell.Screen) {
	if u.isLoading {
		currentText := u.GetText()
		u.SetLabel("URL 🔄 ")
		u.InputField.Draw(screen)
		u.SetLabel("URL: ")
		u.SetText(currentText)
	} else {
		u.InputField.Draw(screen)
	}
}
