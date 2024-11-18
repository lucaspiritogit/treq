package service

import (
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
		requestList.AddItem(request.Title, request.Method, 0, nil)
	}
}
