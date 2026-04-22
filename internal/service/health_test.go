package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockChecker struct {
	err   error
	delay time.Duration
}

func (m *MockChecker) Check(ctx context.Context) error {
	time.Sleep(m.delay)
	return m.err
}

func TestHealthService_Check(t *testing.T) {
	t.Run("All healthy", func(t *testing.T) {
		svc := NewHealthService(&MockChecker{}, &MockChecker{})
		err := svc.Check(context.Background())
		assert.NoError(t, err)
	})

	t.Run("One fails", func(t *testing.T) {
		errFail := errors.New("db down")
		svc := NewHealthService(&MockChecker{err: errFail}, &MockChecker{})
		err := svc.Check(context.Background())
		assert.ErrorContains(t, err, "db down")
	})

	t.Run("Timeout", func(t *testing.T) {
		svc := NewHealthService(&MockChecker{delay: time.Hour})
		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
		defer cancel()

		err := svc.Check(ctx)
		assert.ErrorContains(t, err, "health check timeout")
	})
}
