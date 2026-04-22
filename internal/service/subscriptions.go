package service

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/longwavee/effective-mobile-test/internal/model"
	"golang.org/x/sync/errgroup"
)

type (
	SubscriptionRepository interface {
		Add(ctx context.Context, sub *model.Subscription) error
		FindByID(ctx context.Context, id int64) (model.Subscription, error)
		Update(ctx context.Context, sub *model.Subscription) error
		Remove(ctx context.Context, id int64) error
		ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.Subscription, error)
	}
)

type (
	SubscriptionService struct {
		repo SubscriptionRepository
	}
)

func NewSubscriptionService(repo SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{
		repo: repo,
	}
}

func (s *SubscriptionService) Create(
	ctx context.Context,
	sub *model.Subscription,
) error {
	if err := s.validate(sub); err != nil {
		return err
	}
	return s.repo.Add(ctx, sub)
}

func (s *SubscriptionService) GetByID(
	ctx context.Context,
	id int64,
) (model.Subscription, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *SubscriptionService) Update(
	ctx context.Context,
	sub *model.Subscription,
) error {
	if err := s.validate(sub); err != nil {
		return err
	}
	return s.repo.Update(ctx, sub)
}

func (s *SubscriptionService) Delete(
	ctx context.Context,
	id int64,
) error {
	return s.repo.Remove(ctx, id)
}

func (s *SubscriptionService) ListByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]model.Subscription, error) {
	return s.repo.ListByUserID(ctx, userID)
}

func (s *SubscriptionService) TotalCostForPeriod(
	ctx context.Context,
	userID uuid.UUID,
	serviceName string,
	from time.Time,
	to time.Time,
) (int64, error) {
	subs, err := s.repo.ListByUserID(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("getting subs: %w", err)
	}

	g, ctx := errgroup.WithContext(ctx)

	var cost int64

	for _, sub := range subs {
		g.Go(func() error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			if serviceName != "" && sub.ServiceName != serviceName {
				return nil
			}

			if sub.StartDate.After(to) {
				return nil
			}
			if sub.EndDate != nil && sub.EndDate.Before(from) {
				return nil
			}

			atomic.AddInt64(
				&cost,
				s.calculateSubCost(sub, from, to),
			)

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return 0, err
	}

	return cost, nil
}

func (s *SubscriptionService) calculateSubCost(sub model.Subscription, from, to time.Time) int64 {
	actualStart := sub.StartDate
	if from.After(actualStart) {
		actualStart = from
	}

	actualEnd := to
	if sub.EndDate != nil && sub.EndDate.Before(to) {
		actualEnd = *sub.EndDate
	}

	years := actualEnd.Year() - actualStart.Year()
	months := int(actualEnd.Month()) - int(actualStart.Month())
	totalMonths := years*12 + months + 1

	if totalMonths <= 0 {
		return 0
	}

	return int64(totalMonths) * int64(sub.Price)
}

func (s *SubscriptionService) validate(sub *model.Subscription) error {
	if sub.ServiceName == "" {
		return model.ErrSubscriptionEmptyService
	}

	if sub.EndDate != nil && !sub.EndDate.IsZero() {
		if sub.StartDate.After(*sub.EndDate) {
			return model.ErrSubscriptionInvalidPeriod
		}
	}
	return nil
}
