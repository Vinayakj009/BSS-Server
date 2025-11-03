package database

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCreatePlan(t *testing.T) {
	ctx, db := createDbForPlanTests(t)
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
	ctx, db := createDbForPlanTests(t)
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
	ctx, db := createDbForPlanTests(t)

	// Replace with an existing plan ID in your test database
	existingPlanID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	if plan, err := db.GetPlan(ctx, existingPlanID); err != nil {
		t.Fatalf("Failed to get plan: %v", err)
	} else {
		t.Logf("Successfully retrieved plan: ID=%s, Name=%s", plan.ID, plan.Name)
	}
}

func TestUpdatePlan(t *testing.T) {
	ctx, db := createDbForPlanTests(t)
	defer db.Close()
	// Replace with an existing plan ID in your test database
	existingPlanID := uuid.MustParse("11111111-1111-1111-1111-111111111111")

	// Fetch the existing plan
	plan, err := db.GetPlan(ctx, existingPlanID)
	if err != nil {
		t.Fatalf("Failed to get plan: %v", err)
	}

	// Update plan details
	plan.Name = "Updated Test Plan"
	plan.PriceCents = 1999
	plan.UpdatedAt = time.Now()

	if updatedPlan, err := db.UpdatePlan(ctx, plan); err != nil {
		t.Fatalf("Failed to update plan: %v", err)
	} else {
		t.Logf("Successfully updated plan: ID=%s, Name=%s", updatedPlan.ID, updatedPlan.Name)
	}
}
