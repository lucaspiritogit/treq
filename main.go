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

	requestList := request.NewRequestList()
	requestRepository, err := repository.NewRequestRepository(requestList.TreeView, app)
	if err != nil {
		panic(err)
	}

	requestService := service.NewRequestService(requestList.TreeView, requestRepository)
	app.EnableMouse(true)
	app.EnablePaste(true)
	appFlex := ui.InitializeAppUI(app, requestList.TreeView, requestService, requestRepository)

	if err := app.SetRoot(appFlex, true).Run(); err != nil {
		panic(err)
	}
}
