package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/longwavee/effective-mobile-test/internal/config"
	"github.com/longwavee/effective-mobile-test/internal/model"
)

const (
	defaultMaxConns        = 10
	defaultMinConns        = 2
	defaultMaxConnLifetime = 1 * time.Hour
	defaultMaxConnIdleTime = 30 * time.Minute
)

type (
	SubscriptionRepo struct {
		pool *pgxpool.Pool
	}
)

func NewSubscriptionRepo(
	ctx context.Context,
	cfg *config.Postgres,
) (*SubscriptionRepo, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.ConnString())
	if err != nil {
		return nil, fmt.Errorf("parse config failed: %w", err)
	}

	poolCfg.MaxConns = defaultMaxConns
	poolCfg.MinConns = defaultMinConns
	poolCfg.MaxConnLifetime = defaultMaxConnLifetime
	poolCfg.MaxConnIdleTime = defaultMaxConnIdleTime

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}

	return &SubscriptionRepo{pool: pool}, nil
}

func (r *SubscriptionRepo) Close() {
	if r.pool != nil {
		r.pool.Close()
	}
}

func (r *SubscriptionRepo) Check(ctx context.Context) error {
	if err := r.pool.Ping(ctx); err != nil {
		return fmt.Errorf("ping failed: %w", err)
	}

	return nil
}

func (r *SubscriptionRepo) Add(
	ctx context.Context,
	sub *model.Subscription,
) error {
	query := `
        INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id
    `

	err := r.pool.QueryRow(ctx, query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	).Scan(&sub.ID)

	if err != nil {
		return fmt.Errorf("add subscription failed: %w", err)
	}
	return nil
}

func (r *SubscriptionRepo) FindByID(
	ctx context.Context,
	id int64,
) (model.Subscription, error) {
	query := `
        SELECT id, service_name, price, user_id, start_date, end_date
        FROM subscriptions
        WHERE id = $1
    `

	var sub model.Subscription
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return sub, model.ErrSubscriptionNotFound
		}
		return sub, fmt.Errorf("find subscription by id failed: %w", err)
	}
	return sub, nil
}

func (r *SubscriptionRepo) Update(
	ctx context.Context,
	sub *model.Subscription,
) error {
	query := `
        UPDATE subscriptions
        SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5
        WHERE id = $6
    `

	tag, err := r.pool.Exec(ctx, query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
		sub.ID,
	)
	if err != nil {
		return fmt.Errorf("update subscription failed: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return model.ErrSubscriptionNotFound
	}
	return nil
}

func (r *SubscriptionRepo) Remove(
	ctx context.Context,
	id int64,
) error {
	query := `
        DELETE FROM subscriptions
        WHERE id = $1
    `

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("remove subscription failed: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return model.ErrSubscriptionNotFound
	}
	return nil
}

func (r *SubscriptionRepo) ListByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]model.Subscription, error) {
	query := `
		SELECT id, service_name, price, user_id, start_date, end_date
		FROM subscriptions
		WHERE user_id = $1
	`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("list subscriptions by user id failed: %w", err)
	}
	defer rows.Close()

	subs := []model.Subscription{}
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
			return nil, fmt.Errorf("scan subscription failed: %w", err)
		}
		subs = append(subs, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate subscriptions failed: %w", err)
	}
	return subs, nil
}
