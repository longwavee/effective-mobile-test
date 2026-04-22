package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/longwavee/effective-mobile-test/internal/model"
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

func (s *SubscriptionService) Create(ctx context.Context, sub *model.Subscription) error {
	if err := s.validate(sub); err != nil {
		return fmt.Errorf("service: validate subscription: %w", err)
	}
	return s.repo.Add(ctx, sub)
}

func (s *SubscriptionService) GetByID(ctx context.Context, id int64) (model.Subscription, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *SubscriptionService) Update(ctx context.Context, sub *model.Subscription) error {
	if err := s.validate(sub); err != nil {
		return fmt.Errorf("service: validate subscription: %w", err)
	}
	return s.repo.Update(ctx, sub)
}

func (s *SubscriptionService) Delete(ctx context.Context, id int64) error {
	return s.repo.Remove(ctx, id)
}

func (s *SubscriptionService) ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.Subscription, error) {
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
		return 0, fmt.Errorf("service: get user subscriptions: %w", err)
	}

	var cost int64
	for _, sub := range subs {
		if serviceName != "" && sub.ServiceName != serviceName {
			continue
		}

		if !s.isSubActiveInPeriod(sub, from, to) {
			continue
		}

		cost += s.calculateSubCost(sub, from, to)
	}

	return cost, nil
}

func (s *SubscriptionService) isSubActiveInPeriod(sub model.Subscription, from, to time.Time) bool {
	if sub.StartDate.After(to) {
		return false
	}
	if sub.EndDate != nil && !sub.EndDate.IsZero() && sub.EndDate.Before(from) {
		return false
	}
	return true
}

func (s *SubscriptionService) calculateSubCost(sub model.Subscription, from, to time.Time) int64 {
	start := sub.StartDate
	if from.After(start) {
		start = from
	}

	end := to
	if sub.EndDate != nil && !sub.EndDate.IsZero() && sub.EndDate.Before(to) {
		end = *sub.EndDate
	}

	years := end.Year() - start.Year()
	months := int(end.Month()) - int(start.Month())
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

	if sub.Price < 0 {
		sub.Price = 0
	}

	if sub.EndDate != nil && !sub.EndDate.IsZero() {
		if sub.StartDate.After(*sub.EndDate) {
			return model.ErrSubscriptionInvalidPeriod
		}
	}
	return nil
}
