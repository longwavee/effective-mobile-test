package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

const (
	baseAPIPath = "/api"

	APIPathV1 = baseAPIPath + "/v1"
)

const (
	PathLiveness  = "/liveness"
	PathReadiness = "/readiness"

	PathSubscriptions = "/subscriptions"
)

type (
	HealthHandler interface {
		Live(w http.ResponseWriter, r *http.Request)
		Ready(w http.ResponseWriter, r *http.Request)
	}

	SubscriptionHandler interface {
		Create(w http.ResponseWriter, r *http.Request)
		GetByID(w http.ResponseWriter, r *http.Request)
		Update(w http.ResponseWriter, r *http.Request)
		Delete(w http.ResponseWriter, r *http.Request)
		ListByUserID(w http.ResponseWriter, r *http.Request)
		TotalCostForPeriod(w http.ResponseWriter, r *http.Request)
	}
)

func NewHTTPRouter(
	healthHandler HealthHandler,
	subsHandler SubscriptionHandler,

	loggingMiddleware func(http.Handler) http.Handler,
	recoveryMiddleware func(http.Handler) http.Handler,
) http.Handler {
	r := chi.NewRouter()

	if recoveryMiddleware != nil {
		r.Use(recoveryMiddleware)
	}
	if loggingMiddleware != nil {
		r.Use(loggingMiddleware)
	}

	r.Route(APIPathV1, func(r chi.Router) {
		r.Get(PathLiveness, healthHandler.Live)
		r.Get(PathReadiness, healthHandler.Ready)

		r.Route(PathSubscriptions, func(r chi.Router) {
			r.Post("/", subsHandler.Create)
			r.Get("/{id}", subsHandler.GetByID)
			r.Put("/{id}", subsHandler.Update)
			r.Delete("/{id}", subsHandler.Delete)
			r.Get("/list", subsHandler.ListByUserID)
			r.Get("/total-cost", subsHandler.TotalCostForPeriod)
		})
	})

	return r
}
