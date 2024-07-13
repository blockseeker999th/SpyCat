package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		log = log.With(
			slog.String("component", "middlware/logger"),
		)

		log.Info("Logger middlware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := log.With(
				slog.String("Method", r.Method),
				slog.String("Path", r.URL.Path),
				slog.String("Remote_addr", r.RemoteAddr),
				slog.String("User_agent", r.UserAgent()),
			)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			currentTime := time.Now()

			defer func() {
				entry.Info("request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(currentTime).String()),
				)
			}()

			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
