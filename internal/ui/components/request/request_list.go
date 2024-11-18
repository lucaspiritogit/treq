package request

import "github.com/rivo/tview"

func GetRequestList() *tview.List {
	requestList := tview.NewList()

	requestList.SetBorder(true).SetTitle("Requests")

	return requestList
}
