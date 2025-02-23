package middleware

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

type key string

const TraceIDKey key = "traceID"

// Middleware для генерации trace_id
func TraceMeddleWare(newt http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		traceID := r.Header.Get("X-Request-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}

		ctx := context.WithValue(r.Context(), TraceIDKey, traceID)
		r = r.WithContext(ctx)

		// Добавляем trace_id в заголовок ответа
		w.Header().Set("X-Request-ID", traceID)

		newt.ServeHTTP(w, r)
	})
}

// Функция для получения trace_id из контекста
func GetTraceID(ctx context.Context) string {
	traceID, ok := ctx.Value(TraceIDKey).(string)
	if !ok || traceID == "" {
		return uuid.New().String()
	}

	return traceID
}
