package handlers

import (
	"time"

	"github.com/google/uuid"
	"github.com/longwavee/effective-mobile-test/internal/model"
)

const (
	subscriptionDateFormat = "01-2006"
)

type (
	SubscriptionRequest struct {
		ID          int64   `json:"id"`
		ServiceName string  `json:"service_name"`
		Price       int     `json:"price"` // rubles
		UserID      string  `json:"user_id"`
		StartDate   string  `json:"start_date"`
		EndDate     *string `json:"end_date,omitempty"`
	}
)

func (r *SubscriptionRequest) ToModel() (*model.Subscription, error) {
	uid, err := uuid.Parse(r.UserID)
	if err != nil {
		return nil, err
	}

	start, err := time.Parse(subscriptionDateFormat, r.StartDate)
	if err != nil {
		return nil, err
	}

	var end *time.Time
	if r.EndDate != nil && *r.EndDate != "" {
		t, err := time.Parse(subscriptionDateFormat, *r.EndDate)
		if err != nil {
			return nil, err
		}
		end = &t
	}

	return &model.Subscription{
		ID:          r.ID,
		ServiceName: r.ServiceName,
		Price:       r.Price,
		UserID:      uid,
		StartDate:   start,
		EndDate:     end,
	}, nil
}
