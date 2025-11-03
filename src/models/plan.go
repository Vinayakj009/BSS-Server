package models

import (
	"time"

	"github.com/google/uuid"
)

type Plan struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Code         string    `json:"code" db:"code"`
	Name         string    `json:"name" db:"name"`
	PriceCents   int64     `json:"price_cents" db:"price_cents"`
	Currency     string    `json:"currency" db:"currency"`
	DurationDays int       `json:"duration_days" db:"duration_days"`
	DataMB       int64     `json:"data_mb" db:"data_mb"`
	Active       bool      `json:"active" db:"active"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
