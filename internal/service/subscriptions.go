package service

import (
	"context"

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
