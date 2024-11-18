package main

import (
	"treq/internal/repository"
	"treq/internal/service"
	"treq/internal/ui"
	"treq/internal/ui/components/request"

	_ "github.com/mattn/go-sqlite3"

	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	requestList := request.GetRequestList()

	requestRepository, err := repository.NewRequestRepository(requestList, app)
	if err != nil {
		panic(err)
	}

	requestService := service.NewRequestService(requestList, requestRepository)
	requestService.RefreshRequestsList(requestList, requestRepository)

	appFlex := ui.InitializeAppUI(app, requestList, requestService, requestRepository)

	if err := app.SetRoot(appFlex, true).Run(); err != nil {
		panic(err)
	}
}
