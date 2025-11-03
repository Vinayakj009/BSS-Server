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

func (db *DB) GetSubscriptionsByUserId(ctx context.Context, pageableRequest PageableRequest, userId string) (Page[Subscription], error) {
	offset := (pageableRequest.Page - 1) * pageableRequest.PageSize
	pageSize := pageableRequest.PageSize
	query := `SELECT * 
			  FROM subscriptions 
			  WHERE customer_id = $1
			  ORDER BY created_at DESC
			  LIMIT $2 OFFSET $3`
	rows, err := db.Pool.Query(ctx, query, userId, pageSize, offset)
	if err != nil {
		return Page[Subscription]{}, err
	}
	defer rows.Close()

	var subscriptions []Subscription
	for rows.Next() {
		subscription, err := scanSubscription(rows)
		if err != nil {
			return Page[Subscription]{}, err
		}
		subscriptions = append(subscriptions, subscription)
	}
	var totalCount int64
	countQuery := `SELECT COUNT(*) 
				   FROM subscriptions 
				   WHERE customer_id = $1`
	err = db.Pool.QueryRow(ctx, countQuery, userId).Scan(&totalCount)
	if err != nil {
		return Page[Subscription]{}, err
	}
	return Page[Subscription]{
		TotalCount: totalCount,
		Items:      subscriptions,
	}, nil
}

func (db *DB) GetActiveSubscriptionByUserId(ctx context.Context, userId string) (Subscription, error) {
	query := `SELECT * 
			  FROM subscriptions 
			  WHERE customer_id = $1 AND status = 'ACTIVE'`
	row := db.Pool.QueryRow(ctx, query, userId)
	return scanSubscription(row)
}
