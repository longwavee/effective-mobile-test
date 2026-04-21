package handlers

import (
	"time"

	"github.com/google/uuid"
	"github.com/longwavee/effective-mobile-test/internal/model"
)

const subscriptionRequestDateFormat = "01-2006"

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
	userID, err := uuid.Parse(r.UserID)
	if err != nil {
		return nil, err
	}

	startDate, err := time.Parse(subscriptionRequestDateFormat, r.StartDate)
	if err != nil {
		return nil, err
	}

	var endDate *time.Time
	if r.EndDate != nil {
		newEndDate, err := time.Parse(subscriptionRequestDateFormat, *r.EndDate)
		if err != nil {
			return nil, err
		}
		endDate = &newEndDate
	}

	return &model.Subscription{
		ID:          r.ID,
		ServiceName: r.ServiceName,
		Price:       r.Price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}, nil
}
