package main

import (
	"treq/internal/repository"
	"treq/internal/service"
	"treq/internal/ui"
	"treq/internal/ui/components/request"
	"treq/internal/ui/state"

	_ "modernc.org/sqlite"

	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	requestList := request.GetRequestList()
	appState := state.NewAppState()

	requestRepository, err := repository.NewRequestRepository(requestList, app)
	if err != nil {
		panic(err)
	}

	requestService := service.NewRequestService(requestList, requestRepository)
	requestService.RefreshRequestsList(requestList, requestRepository)

	appFlex := ui.InitializeAppUI(app, requestList, requestService, requestRepository, appState)

	if err := app.SetRoot(appFlex, true).Run(); err != nil {
		panic(err)
	}
}
