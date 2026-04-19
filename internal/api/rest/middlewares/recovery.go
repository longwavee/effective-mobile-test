package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
)

func Recovery(log *slog.Logger) func(http.Handler) http.Handler {
	if log == nil {
		log = slog.Default()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				if rec := recover(); rec != nil {

					log.Error("panic recovered",
						"error", rec,
						"method", r.Method,
						"path", r.URL.Path,
					)

					if ww.Status() == 0 {
						w.WriteHeader(http.StatusInternalServerError)
						_, _ = w.Write([]byte("Internal server error"))
					}
				}
			}()

			next.ServeHTTP(ww, r)
		})
	}
}
