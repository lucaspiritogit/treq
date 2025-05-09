package service

import (
	"treq/internal/repository"

	"github.com/rivo/tview"
)

type IRequestService interface {
	RefreshRequestsList(requestList *tview.TreeView, requestManager *repository.RequestRepository)
}

type RequestService struct {
	requestList       *tview.TreeView
	requestRepository *repository.RequestRepository
}

func NewRequestService(requestList *tview.TreeView, requestManager *repository.RequestRepository) *RequestService {
	return &RequestService{requestList: requestList, requestRepository: requestManager}
}