package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func (db *DB) CreatePlan(ctx context.Context, plan Plan) (Plan, error) {
	query := `INSERT INTO plans (code, name, price_cents, currency, duration_days, data_mb, created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	var id uuid.UUID
	err := db.Pool.QueryRow(ctx, query,
		plan.Code,
		plan.Name,
		plan.PriceCents,
		plan.Currency,
		plan.DurationDays,
		plan.DataMB,
		plan.CreatedAt,
		plan.UpdatedAt,
	).Scan(&id)
	if err != nil {
		return Plan{}, err
	}
	plan.ID = id
	return plan, nil
}

func (db *DB) scanPlan(ctx context.Context, row pgx.Row) (Plan, error) {
	var plan Plan
	err := row.Scan(
		&plan.ID,
		&plan.Code,
		&plan.Name,
		&plan.PriceCents,
		&plan.Currency,
		&plan.DurationDays,
		&plan.DataMB,
		&plan.Active,
		&plan.CreatedAt,
		&plan.UpdatedAt,
	)
	return plan, err
}

func (db *DB) GetPlans(ctx context.Context, pageableRequest PageableRequest) (Page[Plan], error) {
	offset := (pageableRequest.Page - 1) * pageableRequest.PageSize
	pageSize := pageableRequest.PageSize
	query := `SELECT * from plans ORDER BY created_at DESC LIMIT $1 OFFSET $2`
	rows, err := db.Pool.Query(ctx, query, pageSize, offset)
	if err != nil {
		return Page[Plan]{}, err
	}
	defer rows.Close()
	var plans []Plan
	for rows.Next() {
		plan, err := db.scanPlan(ctx, rows)
		if err != nil {
			return Page[Plan]{}, err
		}
		plans = append(plans, plan)
	}
	var totalCount int64
	countQuery := `SELECT COUNT(*) FROM plans`
	err = db.Pool.QueryRow(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		return Page[Plan]{}, err
	}
	return Page[Plan]{
		TotalCount: totalCount,
		Items:      plans,
	}, nil
}

func (db *DB) GetPlan(ctx context.Context, id uuid.UUID) (Plan, error) {
	query := `SELECT * from plans WHERE id = $1`
	row := db.Pool.QueryRow(ctx, query, id)
	return db.scanPlan(ctx, row)
}
