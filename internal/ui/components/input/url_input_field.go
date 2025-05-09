package input

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
	headers          *Headers
	app              *tview.Application
	isLoading        bool
}

func NewURLInputField(responseView *tview.TextView, responseMetadata *tview.TextView, httpVerbDropdown *tview.DropDown, headers *Headers, app *tview.Application) *URLInputField {
	urlField := &URLInputField{
		InputField:       tview.NewInputField(),
		responseView:     responseView,
		responseMetadata: responseMetadata,
		httpVerbDropdown: httpVerbDropdown,
		headers:          headers,
		app:              app,
	}

	urlField.
		SetFieldWidth(70).
		SetAcceptanceFunc(tview.InputFieldMaxLength(1024))

	urlField.SetBorderPadding(1,1,0,0)

	urlField.SetFieldTextColor(tcell.ColorBlack)
	urlField.SetFieldBackgroundColor(tcell.ColorWhite)
	urlField.SetPlaceholder("Enter URL").
		SetPlaceholderTextColor(tcell.ColorWhite)

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
		return
	case tcell.KeyEsc:
		u.SetBorderColor(tcell.ColorWhite)
		return
	}
}

func (u *URLInputField) executeRequest() {
	u.isLoading = true
	u.SetBorderColor(tcell.ColorYellow)
	u.responseMetadata.SetText("[yellow]Loading...")

	_, verb := u.httpVerbDropdown.GetCurrentOption()

	go func() {
		result, err := http_request.FetchUrl(u.GetText(), verb, u.headers.GetHeaders())
		if err != nil {
			u.app.QueueUpdateDraw(func() {
				u.responseView.SetText("[red]Error: [white]Failed to fetch URL")
				u.responseMetadata.SetText(fmt.Sprintf("%s [white]Content-Length: [blue]%d", "-", 0))
				u.SetBorderColor(tcell.ColorRed)
				u.isLoading = false
			})
			return
		}
		u.app.QueueUpdateDraw(func() {
			formattedBody := colorizeJSON(result.Body)
			u.responseView.SetText(formattedBody)
			statusText := u.formatStatusCode(result.Resp.StatusCode)
			contentLength := len(result.Body)
			u.responseMetadata.SetText(fmt.Sprintf("%s [white]Content-Length: [blue]%d", statusText, contentLength))
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
		u.SetLabel("🕛")
		u.InputField.Draw(screen)
		u.SetLabel("🕑")
		u.InputField.Draw(screen)
		u.SetLabel("")
		u.SetText(currentText)
	} else {
		u.InputField.Draw(screen)
	}
}
