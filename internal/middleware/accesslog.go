package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"ProjectMessenger/internal/misc"
)

func AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, "traceID", misc.RandStringRunes(8))
		logger := slog.With("requestID", ctx.Value("traceID"))
		logger.Info("request accessLog", "path", r.URL.Path)
		start := time.Now()
		next.ServeHTTP(w, r.WithContext(ctx))
		logger.Info("requestProcessed", "method", r.Method, "remoteAddr", r.RemoteAddr, "URLPath",
			r.URL.Path, "time", time.Since(start))
	})
}
