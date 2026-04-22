package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubscriptionRequest_ToModel(t *testing.T) {
	endDateStr := "12-2024"

	tests := []struct {
		name    string
		request SubscriptionRequest
		wantErr bool
	}{
		{
			name: "Success - valid data",
			request: SubscriptionRequest{
				UserID:    "550e8400-e29b-41d4-a716-446655440000",
				StartDate: "01-2024",
				Price:     100,
			},
			wantErr: false,
		},
		{
			name: "Success - with end date",
			request: SubscriptionRequest{
				UserID:    "550e8400-e29b-41d4-a716-446655440000",
				StartDate: "01-2024",
				EndDate:   &endDateStr,
			},
			wantErr: false,
		},
		{
			name: "Fail - invalid UUID",
			request: SubscriptionRequest{
				UserID:    "invalid-uuid",
				StartDate: "01-2024",
			},
			wantErr: true,
		},
		{
			name: "Fail - invalid start date format",
			request: SubscriptionRequest{
				UserID:    "550e8400-e29b-41d4-a716-446655440000",
				StartDate: "2024-01-01",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, err := tt.request.ToModel()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, model)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, model)
				assert.Equal(t, tt.request.Price, model.Price)
			}
		})
	}
}
