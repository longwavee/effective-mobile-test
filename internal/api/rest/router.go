package rest

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

const (
	baseAPIPath = "/api"

	APIPathV1 = baseAPIPath + "/v1"
)

func NewHTTPRouter(
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

	r.Route(APIPathV1, func(r chi.Router) {})

	return r
}
