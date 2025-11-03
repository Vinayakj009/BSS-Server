package database

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGetSubscriptionsByUserId(t *testing.T) {

	ctx, db := createDbForPlanTests(t)
	defer db.Close()

	pageableRequest := PageableRequest{
		Page:     1,
		PageSize: 10,
	}

	// Replace with an existing user ID in your test database
	existingUserID := "00000000-0000-0000-0000-000000000000"
	subscriptions, err := db.GetSubscriptionsByUserId(ctx, pageableRequest, existingUserID)
	if err != nil {
		t.Fatalf("Failed to get subscription: %v", err)
	}
	if len(subscriptions.Items) != 2 {
		t.Fatalf("Expected 2 subscriptions, got %d", len(subscriptions.Items))
	}
}

func TestGetActiveSubscriptionByUserId(t *testing.T) {
	ctx, db := createDbForPlanTests(t)
	defer db.Close()

	existingUserID := "00000000-0000-0000-0000-000000000000"
	_, err := db.GetActiveSubscriptionByUserId(ctx, existingUserID)
	if err != nil {
		t.Fatalf("Failed to get subscription: %v", err)
	}

	expirtedUserId := "00000000-0000-0000-0000-000000000001"
	emptySubscription, err := db.GetActiveSubscriptionByUserId(ctx, expirtedUserId)
	if err == nil {
		t.Fatalf("Expected no subscription found, but got: %v", emptySubscription)
	}

	nonExistingUserId := "00000000-0000-0000-0000-000000000002"
	nonExistingSubscription, err := db.GetActiveSubscriptionByUserId(ctx, nonExistingUserId)
	if err == nil {
		t.Fatalf("Expected no subscription found, but got: %v", nonExistingSubscription)
	}
	t.Logf("Successfully retrieved subscription: UserID=%s", existingUserID)
}

func TestCreateSubscription(t *testing.T) {
	ctx, db := createDbForPlanTests(t)
	defer db.Close()
	subscription, err := db.CreateSubscription(ctx, Subscription{
		CustomerID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
		PlanID:     uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		StartDate:  time.Now(),
		EndDate:    time.Now().Add(30 * 24 * time.Hour),
		Status:     "ACTIVE",
		AutoRenew:  true,
	})
	if err != nil {
		t.Fatalf("Failed to create subscription: %v", err)
	}
	t.Logf("Successfully created subscription with ID: %s", subscription.ID)
}
func TestCreateSubscriptionBadPlan(t *testing.T) {
	ctx, db := createDbForPlanTests(t)
	defer db.Close()
	subscription, err := db.CreateSubscription(ctx, Subscription{
		CustomerID: uuid.MustParse("00000000-0000-0000-0000-000000000003"),
		PlanID:     uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		StartDate:  time.Now(),
		EndDate:    time.Now().Add(30 * 24 * time.Hour),
		Status:     "ACTIVE",
		AutoRenew:  true,
	})
	if err == nil {
		t.Fatalf("Expected failure when creating subscription with bad plan ID, but got success: %v", subscription)
	}
}
