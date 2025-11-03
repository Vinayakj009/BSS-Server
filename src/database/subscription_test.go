package database

import (
	"testing"
)

func TestGetSubscriptionByUserId(t *testing.T) {
	ctx, db := createDbForPlanTests(t)
	defer db.Close()

	// Replace with an existing user ID in your test database
	existingUserID := "00000000-0000-0000-0000-000000000000"
	subscription, err := db.GetSubscriptionByUserId(ctx, existingUserID)
	if err != nil {
		t.Fatalf("Failed to get subscription: %v", err)
	}
	t.Logf("Successfully retrieved subscription: ID=%s, UserID=%s", subscription.ID, subscription.CustomerID)
}

func TestGetActiveSubscriptionByUserId(t *testing.T) {
	ctx, db := createDbForPlanTests(t)
	defer db.Close()

	existingUserID := "00000000-0000-0000-0000-000000000001"
	_, err := db.GetSubscriptionByUserId(ctx, existingUserID)
	if err != nil {
		t.Fatalf("Failed to get subscription: %v", err)
	}

	emptySubscription, err := db.GetActiveSubscriptionByUserId(ctx, existingUserID)
	if err == nil {
		t.Fatalf("Expected no subscription found, but got: %v", emptySubscription)
	}
	t.Logf("Successfully retrieved subscription: UserID=%s", existingUserID)
}
