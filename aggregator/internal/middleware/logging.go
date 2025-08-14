package middleware

import (
	"aggregator/internal/common/logger"
	"net/http"
	"time"
)

type LoggingMiddleware struct {
	logger logger.Logger
}

func NewLoggingMiddleware(logger logger.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		logger: logger,
	}
}

func (lm *LoggingMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		lm.logger.Info(r.Context(),
			"Method: "+r.Method+", Path: "+r.URL.Path+
				", Duration: "+duration.String())
	})
}
