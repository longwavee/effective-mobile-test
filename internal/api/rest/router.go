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
)

type (
	HealthHandler interface {
		Live(w http.ResponseWriter, r *http.Request)
		Ready(w http.ResponseWriter, r *http.Request)
	}
)

func NewHTTPRouter(
	healthHandler HealthHandler,

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
	})

	return r
}
