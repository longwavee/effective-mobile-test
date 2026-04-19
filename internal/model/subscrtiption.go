package model

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrSubscriptionNotFound      = errors.New("subscription not found")
	ErrSubscriptionEmptyService  = errors.New("subscription service name cannot be empty")
	ErrSubscriptionInvalidPeriod = errors.New("subscription start date cannot be after end date")
)

type (
	Subscription struct {
		ID          int64      `json:"id"`
		ServiceName string     `json:"service_name"`
		Price       int        `json:"price"` // rubles
		UserID      uuid.UUID  `json:"user_id"`
		StartDate   time.Time  `json:"start_date"`
		EndDate     *time.Time `json:"end_date,omitempty"`
	}
)
