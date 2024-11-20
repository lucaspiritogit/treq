package models

import "time"

type SavedHeaders struct {
	ID        int       `json:"id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	RequestId int       `json:"request_id"`
	CreatedAt time.Time `json:"created_at"`
}
