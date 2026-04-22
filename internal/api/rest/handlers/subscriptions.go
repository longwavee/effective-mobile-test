package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

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
	}
)

type (
	SubscriptionHandler struct {
		service SubscriptionService
		log     *slog.Logger
	}
)

func NewSubscriptionHandler(service SubscriptionService, log *slog.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{service: service, log: log}
}

func (h *SubscriptionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var subReq SubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&subReq); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	sub, err := subReq.ToModel()
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
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

	var subReq SubscriptionRequest
	if err := json.NewDecoder(r.Body).Decode(&subReq); err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}
	subReq.ID = id

	sub, err := subReq.ToModel()
	if err != nil {
		h.respondError(w, http.StatusBadRequest, "Invalid request body")
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
		h.respondError(w, http.StatusBadRequest, "Invalid user id")
		return
	}

	subs, err := h.service.ListByUserID(r.Context(), userID)
	if err != nil {
		h.handleServiceError(w, err)
		return
	}

	h.respond(w, http.StatusOK, subs)
}

func (h *SubscriptionHandler) parseID(r *http.Request) (int64, error) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid resource identity")
	}
	return id, nil
}

func (h *SubscriptionHandler) respond(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}

func (h *SubscriptionHandler) respondError(w http.ResponseWriter, code int, msg string) {
	h.respond(w, code, map[string]string{"error": msg})
}

func (h *SubscriptionHandler) handleServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, model.ErrSubscriptionNotFound):
		h.respondError(w, http.StatusNotFound, "Subscription not found")
	case errors.Is(err, model.ErrSubscriptionEmptyService):
		h.respondError(w, http.StatusBadRequest, "Service name is empty")
	case errors.Is(err, model.ErrSubscriptionInvalidPeriod):
		h.respondError(w, http.StatusBadRequest, "Invalid period")

	default:
		h.log.Error("service error", "error", err)
		h.respondError(w, http.StatusInternalServerError, "Internal server error")
	}
}
