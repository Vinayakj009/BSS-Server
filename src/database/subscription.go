package database

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func scanSubscription(row pgx.Row) (Subscription, error) {
	var subscription Subscription
	err := row.Scan(&subscription.ID,
		&subscription.CustomerID,
		&subscription.PlanID,
		&subscription.StartDate,
		&subscription.EndDate,
		&subscription.Status,
		&subscription.AutoRenew,
		&subscription.CreatedAt,
		&subscription.UpdatedAt)
	return subscription, err
}

func (db *DB) GetSubscriptionByUserId(ctx context.Context, userId string) (Subscription, error) {
	query := `SELECT * 
			  FROM subscriptions 
			  WHERE customer_id = $1`
	row := db.Pool.QueryRow(ctx, query, userId)
	return scanSubscription(row)
}
