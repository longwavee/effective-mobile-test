package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/longwavee/effective-mobile-test/internal/model"
)

type (
	Querier interface {
		Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
		Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
		QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	}
)

type (
	SubscriptionRepo struct {
		db Querier
	}
)

func NewSubscriptionRepo(db Querier) *SubscriptionRepo {
	return &SubscriptionRepo{db: db}
}

func (r *SubscriptionRepo) Add(ctx context.Context, sub *model.Subscription) error {
	const query = `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	err := r.db.QueryRow(ctx, query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	).Scan(&sub.ID)

	if err != nil {
		return fmt.Errorf("repository: create subscription: %w", err)
	}
	return nil
}

func (r *SubscriptionRepo) FindByID(ctx context.Context, id int64) (model.Subscription, error) {
	const query = `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE id = $1
	`
	var sub model.Subscription
	err := r.db.QueryRow(ctx, query, id).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Subscription{}, model.ErrSubscriptionNotFound
		}
		return model.Subscription{}, fmt.Errorf("repository: get subscription by id: %w", err)
	}
	return sub, nil
}

func (r *SubscriptionRepo) Update(ctx context.Context, sub *model.Subscription) error {
	const query = `
		UPDATE subscriptions
		SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5
		WHERE id = $6
	`
	tag, err := r.db.Exec(ctx, query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
		sub.ID,
	)
	if err != nil {
		return fmt.Errorf("repository: update subscription: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return model.ErrSubscriptionNotFound
	}
	return nil
}

func (r *SubscriptionRepo) Remove(ctx context.Context, id int64) error {
	const query = `
		DELETE FROM subscriptions
		WHERE id = $1
	`
	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("repository: delete subscription: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return model.ErrSubscriptionNotFound
	}
	return nil
}

func (r *SubscriptionRepo) ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.Subscription, error) {
	const query = `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE user_id = $1
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("repository: list subscriptions: %w", err)
	}
	defer rows.Close()

	subs := make([]model.Subscription, 0)

	for rows.Next() {
		var sub model.Subscription
		if err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate,
		); err != nil {
			return nil, fmt.Errorf("repository: scan subscription: %w", err)
		}
		subs = append(subs, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repository: rows error: %w", err)
	}
	return subs, nil
}
