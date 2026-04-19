package service

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

type (
	HealthChecker interface {
		Check(context.Context) error
	}
)

type (
	HealthService struct {
		checkers []HealthChecker
	}
)

func NewHealthService(checkers ...HealthChecker) *HealthService {
	return &HealthService{
		checkers: checkers,
	}
}

func (u *HealthService) Check(ctx context.Context) error {
	var (
		wg     sync.WaitGroup
		mu     sync.Mutex
		resErr error
	)

	for _, checker := range u.checkers {
		wg.Go(func() {
			if err := checker.Check(ctx); err != nil {
				mu.Lock()
				resErr = errors.Join(resErr, err)
				mu.Unlock()
			}
		})
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return resErr
	case <-ctx.Done():
		return fmt.Errorf("health check timeout: %w", ctx.Err())
	}
}
