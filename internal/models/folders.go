package models

import "time"

type Folders struct {
  ID        int    `json:"id"`
  Name      string `json:"name"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
  DeletedAt *time.Time `json:"deleted_at"`
  Requests  []SavedRequest `json:"requests"`
}