package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/longwavee/effective-mobile-test/internal/config"
)

type (
	HealthProvider interface {
		Check(ctx context.Context) error
	}
)

type (
	HealthHandler struct {
		provider         HealthProvider
		readinessTimeout time.Duration
	}
)

func NewHealthHandler(cfg *config.HTTPServer, provider HealthProvider) *HealthHandler {
	return &HealthHandler{
		provider:         provider,
		readinessTimeout: cfg.ReadinessTimeout,
	}
}

func (h *HealthHandler) Live(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte("OK"))
}

func (h *HealthHandler) Ready(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	ctx, cancel := context.WithTimeout(r.Context(), h.readinessTimeout)
	defer cancel()

	err := h.provider.Check(ctx)

	if err != nil {
		// TODO: add a normal writer
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("Service unavailable"))
	} else {
		// TODO: add a normal writer
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}
}
