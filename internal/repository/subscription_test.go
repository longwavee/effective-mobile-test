package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/longwavee/effective-mobile-test/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	connString := "postgres://postgres:password@localhost:5432/database?sslmode=disable"
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		t.Skip("Postgres not available for integration tests")
	}
	return pool
}

func TestSubscriptionRepo_CRUD(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()
	repo := NewSubscriptionRepo(db)
	ctx := context.Background()

	userID := uuid.New()
	sub := &model.Subscription{
		ServiceName: "Netflix",
		Price:       990,
		UserID:      userID,
		StartDate:   time.Now().Truncate(time.Second),
	}

	t.Run("Create and Find", func(t *testing.T) {
		err := repo.Add(ctx, sub)
		require.NoError(t, err)
		assert.NotZero(t, sub.ID)

		found, err := repo.FindByID(ctx, sub.ID)
		require.NoError(t, err)
		assert.Equal(t, sub.ServiceName, found.ServiceName)
		assert.Equal(t, sub.UserID, found.UserID)
	})

	t.Run("Update", func(t *testing.T) {
		sub.Price = 1200
		err := repo.Update(ctx, sub)
		assert.NoError(t, err)

		found, _ := repo.FindByID(ctx, sub.ID)
		assert.Equal(t, 1200, found.Price)
	})

	t.Run("Find Non-Existent", func(t *testing.T) {
		_, err := repo.FindByID(ctx, 999999)
		assert.ErrorIs(t, err, model.ErrSubscriptionNotFound)
	})

	t.Run("Remove", func(t *testing.T) {
		err := repo.Remove(ctx, sub.ID)
		assert.NoError(t, err)

		_, err = repo.FindByID(ctx, sub.ID)
		assert.ErrorIs(t, err, model.ErrSubscriptionNotFound)
	})
}
