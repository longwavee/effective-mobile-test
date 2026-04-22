package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/longwavee/effective-mobile-test/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Add(ctx context.Context, sub *model.Subscription) error {
	return m.Called(ctx, sub).Error(0)
}
func (m *MockRepository) FindByID(ctx context.Context, id int64) (model.Subscription, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(model.Subscription), args.Error(1)
}
func (m *MockRepository) Update(ctx context.Context, sub *model.Subscription) error {
	return m.Called(ctx, sub).Error(0)
}
func (m *MockRepository) Remove(ctx context.Context, id int64) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.Subscription, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]model.Subscription), args.Error(1)
}

func TestSubscriptionService_Validate(t *testing.T) {
	svc := NewSubscriptionService(nil)

	t.Run("Empty service name", func(t *testing.T) {
		err := svc.validate(&model.Subscription{ServiceName: ""})
		assert.ErrorIs(t, err, model.ErrSubscriptionEmptyService)
	})

	t.Run("Invalid period", func(t *testing.T) {
		start := time.Now()
		end := start.Add(-time.Hour)
		err := svc.validate(&model.Subscription{
			ServiceName: "Netflix",
			StartDate:   start,
			EndDate:     &end,
		})
		assert.ErrorIs(t, err, model.ErrSubscriptionInvalidPeriod)
	})

	t.Run("Negative price reset", func(t *testing.T) {
		sub := &model.Subscription{ServiceName: "Spotify", Price: -100}
		_ = svc.validate(sub)
		assert.Equal(t, 0, sub.Price)
	})
}

func TestSubscriptionService_CalculateSubCost(t *testing.T) {
	svc := NewSubscriptionService(nil)
	price := 500

	tests := []struct {
		name     string
		subStart time.Time
		subEnd   *time.Time
		reqFrom  time.Time
		reqTo    time.Time
		want     int64
	}{
		{
			name:     "Full month sub in full month request",
			subStart: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			reqFrom:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			reqTo:    time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC),
			want:     500, // 1 month
		},
		{
			name:     "Two months crossing year",
			subStart: time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC),
			reqFrom:  time.Date(2023, 12, 1, 0, 0, 0, 0, time.UTC),
			reqTo:    time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			want:     1000, // Dec + Jan
		},
		{
			name:     "Sub ends before request period",
			subStart: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			subEnd:   ptr(time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)),
			reqFrom:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			reqTo:    time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
			want:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sub := model.Subscription{
				StartDate: tt.subStart,
				EndDate:   tt.subEnd,
				Price:     price,
			}
			got := svc.calculateSubCost(sub, tt.reqFrom, tt.reqTo)
			assert.Equal(t, tt.want, got)
		})
	}
}

func ptr(t time.Time) *time.Time { return &t }
