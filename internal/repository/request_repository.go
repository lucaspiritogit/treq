package repository

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"treq/internal/models"

	"github.com/rivo/tview"
	_ "modernc.org/sqlite"
)

type RequestRepository struct {
	DbFilePath              string
	db                      *sql.DB
	list                    *tview.List
	app                     *tview.Application
	RequestListMapIndexToId map[int]int
}

type IRequestRepository interface {
	SaveRequest(request models.SavedRequest)
	GetRequests() []models.SavedRequest
	GetRequestById(id int) []models.SavedRequest
	DeleteRequestById(id int) error
}

func NewRequestRepository(list *tview.List, app *tview.Application) (*RequestRepository, error) {
	appData, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("could not find user config directory: %v", err)
	}

	dbPath := appData + "/treq/treq.db"
	dbDir := filepath.Dir(dbPath)

	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return nil, fmt.Errorf("could not create directory for DB: %v", err)
		}
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	createSavedRequestsTable(db)
	createSavedHeadersTable(db)

	return &RequestRepository{db: db, DbFilePath: dbPath, list: list, app: app}, nil
}

func (r *RequestRepository) DeleteRequestById(id int) error {
	itemId, exists := r.RequestListMapIndexToId[id]
	if !exists {
		return fmt.Errorf("no request found for index %d", itemId)
	}
	_, err := r.db.Exec("DELETE FROM saved_requests WHERE id = ?", itemId)
	if err != nil {
		return fmt.Errorf("could not delete request: %v", err)
	}
	r.list.RemoveItem(itemId)
	return nil
}

func (r *RequestRepository) GetRequestById(id int) models.SavedRequest {
	var request models.SavedRequest
	itemId := r.RequestListMapIndexToId[id]
	row := r.db.QueryRow("SELECT id, title, method, url, body FROM saved_requests WHERE id = ?", itemId)

	err := row.Scan(&request.ID, &request.Title, &request.Method, &request.URL, &request.Body)
	if err != nil {
		fmt.Println(err)
	}

	return request
}

func (r *RequestRepository) GetRequests() []models.SavedRequest {
	var requests []models.SavedRequest

	rows, err := r.db.Query("SELECT id, title, method, url, body FROM saved_requests")
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		var request models.SavedRequest
		err = rows.Scan(&request.ID, &request.Title, &request.Method, &request.URL, &request.Body)
		if err != nil {
			fmt.Println(err)
		}
		requests = append(requests, request)
	}
	return requests
}

func (r *RequestRepository) SaveRequest(request models.SavedRequest) error {
	_, err := r.db.Exec("INSERT INTO saved_requests (method, title, url, body) VALUES (?, ?, ?, ?)", request.Method, request.Title, request.URL, request.Body)
	if err != nil {
		return fmt.Errorf("could not insert request: %v", err)
	}
	savedRequest := models.SavedRequest{
		Method: request.Method,
		URL:    request.URL,
	}

	r.list.AddItem(savedRequest.URL, savedRequest.Method, 0, nil)
	return nil
}

func createSavedRequestsTable(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS saved_requests (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		method TEXT NOT NULL,
		url TEXT NOT NULL,
		body TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		modified_at DATETIME
	)`)
	if err != nil {
		return fmt.Errorf("could not create table: %v", err)
	}

	return nil
}

func createSavedHeadersTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS request_headers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    request_id INTEGER NOT NULL,
    header_key TEXT NOT NULL,
    header_value TEXT NOT NULL,
    FOREIGN KEY (request_id) REFERENCES saved_requests(id) ON DELETE CASCADE
);`)

	if err != nil {
		return fmt.Errorf("could not create table: %v", err)
	}

	return nil
}
