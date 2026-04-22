package handlers

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/longwavee/effective-mobile-test/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestSubscriptionHandler_ParseID(t *testing.T) {
	h := &SubscriptionHandler{}

	tests := []struct {
		name    string
		param   string
		wantID  int64
		wantErr bool
	}{
		{"Valid ID", "10", 10, false},
		{"Zero ID", "0", 0, true},
		{"Negative ID", "-5", 0, true},
		{"Not a number", "abc", 0, true},
		{"Empty", "", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.param)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			id, err := h.parseID(req)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantID, id)
			}
		})
	}
}

func TestSubscriptionHandler_HandleServiceError(t *testing.T) {
	discardLog := slog.New(slog.NewTextHandler(io.Discard, nil))

	h := &SubscriptionHandler{
		baseHandler: baseHandler{log: discardLog},
	}

	tests := []struct {
		name       string
		err        error
		wantStatus int
	}{
		{"Not Found", model.ErrSubscriptionNotFound, http.StatusNotFound},
		{"Empty Service", model.ErrSubscriptionEmptyService, http.StatusBadRequest},
		{"Invalid Period", model.ErrSubscriptionInvalidPeriod, http.StatusBadRequest},
		{"Unknown Error", errors.New("database boom"), http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			h.handleServiceError(w, tt.err)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		})
	}
}
