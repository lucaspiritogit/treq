package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"

	"github.com/rivo/tview"
)

type RequestManager struct {
	DbFilePath              string
	db                      *sql.DB
	list                    *tview.List
	app                     *tview.Application
	requestListMapIndexToId map[int]int
}

type IRequestManager interface {
	SaveRequest(request SavedRequest)
	GetRequests() []SavedRequest
}

func NewRequestManager(list *tview.List, app *tview.Application) (*RequestManager, error) {
	appData, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("could not find user config directory: %v", err)
	}

	dbPath := appData + "/turl/turl.db"
	dbDir := filepath.Dir(dbPath)

	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return nil, fmt.Errorf("could not create directory for DB: %v", err)
		}
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	createSavedRequestsTable(db)

	return &RequestManager{db: db, DbFilePath: dbPath, list: list, app: app}, nil
}

func (r *RequestManager) DeleteRequestById(id int) error {
	itemMapId, exists := r.requestListMapIndexToId[id]
	if !exists {
		return fmt.Errorf("no request found for index %d", itemMapId)
	}
	_, err := r.db.Exec("DELETE FROM saved_requests WHERE id = ?", itemMapId)
	if err != nil {
		return fmt.Errorf("could not delete request: %v", err)
	}
	r.list.RemoveItem(itemMapId)
	return nil
}

func (r *RequestManager) GetRequests() []SavedRequest {
	var requests []SavedRequest

	rows, err := r.db.Query("SELECT id, method, url, body FROM saved_requests")
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		var request SavedRequest
		err = rows.Scan(&request.ID, &request.Method, &request.URL, &request.Body)
		if err != nil {
			fmt.Println(err)
		}
		requests = append(requests, request)
	}
	return requests
}

func RefreshRequestsList(requestList *tview.List, requestManager *RequestManager) {
	savedRequests := requestManager.GetRequests()
	requestList.Clear()
	requestManager.requestListMapIndexToId = make(map[int]int)

	for i, request := range savedRequests {
		requestManager.requestListMapIndexToId[i] = request.ID
		requestList.AddItem(request.URL, request.Method, 0, nil)
	}

}

func (r *RequestManager) SaveRequest(request SavedRequest) error {
	_, err := r.db.Exec("INSERT INTO saved_requests (method, url, body) VALUES (?, ?, ?)", request.Method, request.URL, request.Body)
	if err != nil {
		return fmt.Errorf("could not insert request: %v", err)
	}
	savedRequest := SavedRequest{
		Method: request.Method,
		URL:    request.URL,
	}

	r.list.AddItem(savedRequest.URL, savedRequest.Method, 0, nil)
	return nil
}

func createSavedRequestsTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS saved_requests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		method TEXT NOT NULL,
		url TEXT NOT NULL,
		body TEXT,
		headers TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return fmt.Errorf("could not create table: %v", err)
	}

	return nil
}
