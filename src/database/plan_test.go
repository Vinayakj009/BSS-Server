package database

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCreatePlan(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()

	// Create database connection
	db, err := NewDb(ctx)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create a new plan
	plan := Plan{
		ID:         uuid.New(),
		Name:       "Test Plan",
		PriceCents: 999,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if p, err := db.CreatePlan(ctx, plan); err != nil {
		t.Fatalf("Failed to create plan: %v", err)
	} else {
		t.Logf("Successfully created plan with ID: %s", p.ID)
	}
}

func TestGetPlans(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()
	// Create database connection
	db, err := NewDb(ctx)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	pageableRequest := PageableRequest{
		Page:     1,
		PageSize: 10,
	}
	if page, err := db.GetPlans(ctx, pageableRequest); err != nil {
		t.Fatalf("Failed to get plans: %v", err)
	} else {
		t.Logf("Successfully retrieved %d plans", len(page.Items))
		for _, plan := range page.Items {
			t.Logf("Plan ID: %s, Name: %s", plan.ID, plan.Name)
		}
	}
}

func TestGetPlan(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()
	// Create database connection
	db, err := NewDb(ctx)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Replace with an existing plan ID in your test database
	existingPlanID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	if plan, err := db.GetPlan(ctx, existingPlanID); err != nil {
		t.Fatalf("Failed to get plan: %v", err)
	} else {
		t.Logf("Successfully retrieved plan: ID=%s, Name=%s", plan.ID, plan.Name)
	}
}
