package main

import (
	"treq/internal/repository"
	"treq/internal/service"
	"treq/internal/ui"
	"treq/internal/ui/components/request"

	_ "modernc.org/sqlite"

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
	app.EnableMouse(true)
	app.EnablePaste(true)
	appFlex := ui.InitializeAppUI(app, requestList, requestService, requestRepository)

	if err := app.SetRoot(appFlex, true).Run(); err != nil {
		panic(err)
	}
}
