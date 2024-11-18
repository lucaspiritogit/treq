package service

import (
	"fmt"
	"treq/internal/repository"

	"github.com/rivo/tview"
)

type IRequestService interface {
	RefreshRequestsList(requestList *tview.List, requestManager *repository.RequestRepository)
}

type RequestService struct {
	requestList       *tview.List
	requestRepository *repository.RequestRepository
}

func NewRequestService(requestList *tview.List, requestManager *repository.RequestRepository) *RequestService {
	return &RequestService{requestList: requestList, requestRepository: requestManager}
}

func (s *RequestService) RefreshRequestsList(requestList *tview.List, requestRepository *repository.RequestRepository) {
	savedRequests := requestRepository.GetRequests()
	requestList.Clear()
	requestRepository.RequestListMapIndexToId = make(map[int]int)

	for i, request := range savedRequests {
		requestRepository.RequestListMapIndexToId[i] = request.ID
		var httpMethodColor string
		if request.Method == "GET" {
			httpMethodColor = "[green]"
		} else if request.Method == "POST" {
			httpMethodColor = "[yellow]"
		} else if request.Method == "PUT" {
			httpMethodColor = "[blue]"
		} else if request.Method == "DELETE" {
			httpMethodColor = "[red]"
		}

		formatItem := fmt.Sprintf("%s%s[white] %s", httpMethodColor, request.Method, request.Title)
		requestList.AddItem(formatItem, "", 0, nil)
	}
}
