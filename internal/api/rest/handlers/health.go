package handlers

import (
	"context"
	"log/slog"
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
		baseHandler
		provider         HealthProvider
		readinessTimeout time.Duration
	}
)

func NewHealthHandler(cfg *config.HTTPServer, provider HealthProvider, log *slog.Logger) *HealthHandler {
	return &HealthHandler{
		baseHandler:      baseHandler{log: log},
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
		h.respondError(w, http.StatusServiceUnavailable, "service unavailable")
	} else {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}
}
