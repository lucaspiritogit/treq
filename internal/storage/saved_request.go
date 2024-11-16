package storage

import "time"

type SavedRequest struct {
	ID          int       `json:"id"`
	Method      string    `json:"method"`
	URL         string    `json:"url"`
	Body        string    `json:"body"`
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description"`
}
