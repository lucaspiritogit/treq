package request

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type RequestList struct {
	TreeView *tview.TreeView
}

func NewRequestList() *RequestList {
	requestListTreeView := tview.NewTreeView()
	requestListTreeView.SetRoot(tview.NewTreeNode("Requests").SetColor(tcell.ColorWhite))
	requestListTreeView.SetBorder(true)
	requestList := &RequestList{TreeView: requestListTreeView}
	requestList.setInputCapture()
	return requestList

}

func (r *RequestList) setInputCapture() {
	r.TreeView.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			currentNode := r.TreeView.GetRoot()
			if currentNode != nil {
				r.TreeView.SetCurrentNode(currentNode)
			}
			return nil
		}
		switch event.Rune() {
			case 'n' :
				r.AddFolder("New Folder")
				return nil
			case 'x' :
				selectedNode := r.TreeView.GetCurrentNode()
				if selectedNode != nil {
					r.AddRequestToFolder(selectedNode, "New Request")
				}
				return nil
		}
		return event
	})
}

func (r *RequestList) AddFolder(name string) {
	root := r.TreeView.GetRoot()
	folderNode := tview.NewTreeNode(name).SetColor(tcell.ColorRosyBrown)
	root.AddChild(folderNode)
	r.TreeView.SetCurrentNode(folderNode)
}

func (r *RequestList) AddRequestToFolder(folder *tview.TreeNode, requestName string) {
	requestNode := tview.NewTreeNode(requestName).SetColor(tcell.ColorWhite)
	folder.AddChild(requestNode)
}
