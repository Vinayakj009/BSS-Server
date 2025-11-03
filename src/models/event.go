package models

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID         int64     `json:"id" db:"id"`
	EventType  string    `json:"event_type" db:"event_type"`
	ResourceID uuid.UUID `json:"resource_id" db:"resource_id"`
	Payload    []byte    `json:"payload" db:"payload"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
