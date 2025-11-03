package models

import (
	"time"

	"github.com/google/uuid"
)

type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "ACTIVE"
	SubscriptionStatusCancelled SubscriptionStatus = "CANCELLED"
	SubscriptionStatusExpired   SubscriptionStatus = "EXPIRED"
)

type Subscription struct {
	ID         uuid.UUID          `json:"id" db:"id"`
	CustomerID uuid.UUID          `json:"customer_id" db:"customer_id"`
	PlanID     uuid.UUID          `json:"plan_id" db:"plan_id"`
	StartDate  time.Time          `json:"start_date" db:"start_date"`
	EndDate    time.Time          `json:"end_date" db:"end_date"`
	Status     SubscriptionStatus `json:"status" db:"status"`
	AutoRenew  bool               `json:"auto_renew" db:"auto_renew"`
	CreatedAt  time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at" db:"updated_at"`
}
