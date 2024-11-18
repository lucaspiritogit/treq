package models

import "time"

type SavedRequest struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	Method     string    `json:"method"`
	URL        string    `json:"url"`
	Body       string    `json:"body"`
	CreatedAt  time.Time `json:"created_at"`
	ModifiedAt time.Time `json:"modified_at"`
}
