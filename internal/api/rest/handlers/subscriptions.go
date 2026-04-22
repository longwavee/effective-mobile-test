package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/longwavee/effective-mobile-test/internal/model"
)

type (
	SubscriptionService interface {
		Create(ctx context.Context, sub *model.Subscription) error
		GetByID(ctx context.Context, id int64) (model.Subscription, error)
		Update(ctx context.Context, sub *model.Subscription) error
		Delete(ctx context.Context, id int64) error
		ListByUserID(ctx context.Context, userID uuid.UUID) ([]model.Subscription, error)
		TotalCostForPeriod(ctx context.Context, userID uuid.UUID, serviceName string, periodStart, periodEnd time.Time) (int64, error)
	}
)

type (
	SubscriptionHandler struct {
		baseHandler
		service SubscriptionService
	}
)

func NewSubscriptionHandler(service SubscriptionService, log *slog.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{
		baseHandler: baseHandler{log: log},
		service:     service,
	}
}

func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req SubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	sub, err := req.ToModel()
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid data format")
		return
	}

	if err := h.service.Create(r.Context(), sub); err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.respond(w, http.StatusCreated, sub)
}

func (h *SubscriptionHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseID(r)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	sub, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.respond(w, http.StatusOK, sub)
}

func (h *SubscriptionHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseID(r)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req SubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req.ID = id

	sub, err := req.ToModel()
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.service.Update(r.Context(), sub); err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.respond(w, http.StatusNoContent, nil)
}

func (h *SubscriptionHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseID(r)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.respond(w, http.StatusNoContent, nil)
}

func (h *SubscriptionHandler) ListByUserID(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	subs, err := h.service.ListByUserID(r.Context(), userID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.respond(w, http.StatusOK, subs)
}

func (h *SubscriptionHandler) TotalCostForPeriod(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	userID, err := uuid.Parse(q.Get("user_id"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	start, err := time.Parse(subscriptionDateFormat, q.Get("period_start"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid period_start")
		return
	}

	end, err := time.Parse(subscriptionDateFormat, q.Get("period_end"))
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "invalid period_end")
		return
	}

	cost, err := h.service.TotalCostForPeriod(
		r.Context(),
		userID,
		q.Get("service_name"),
		start,
		end,
	)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.respond(w, http.StatusOK, map[string]int64{"total_cost": cost})
}

func (h *SubscriptionHandler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, model.ErrSubscriptionNotFound):
		h.respondError(w, http.StatusNotFound, model.ErrSubscriptionNotFound.Error())
	case errors.Is(err, model.ErrSubscriptionEmptyService):
		h.respondError(w, http.StatusBadRequest, model.ErrSubscriptionEmptyService.Error())
	case errors.Is(err, model.ErrSubscriptionInvalidPeriod):
		h.respondError(w, http.StatusBadRequest, model.ErrSubscriptionInvalidPeriod.Error())
	default:
		h.log.Error("internal error", "err", err)
		h.respondError(w, http.StatusInternalServerError, "internal server error")
	}
}

func (h *SubscriptionHandler) parseID(r *http.Request) (int64, error) {
	param := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(param, 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid resource identity")
	}
	return id, nil
}
