package middlewares

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func Logging(log *slog.Logger) func(http.Handler) http.Handler {
	if log == nil {
		log = slog.Default()
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				duration := time.Since(start)

				attrs := []slog.Attr{
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.Int("status", ww.Status()),
					slog.Duration("duration", duration),
				}

				if ww.Status() >= http.StatusBadRequest {
					log.LogAttrs(r.Context(), slog.LevelError, "http request failed", attrs...)
				} else {
					log.LogAttrs(r.Context(), slog.LevelInfo, "http request completed", attrs...)
				}
			}()

			next.ServeHTTP(ww, r)
		})
	}
}
