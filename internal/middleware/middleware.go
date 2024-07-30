package middleware

import (
	"fmt"
	"net/http"

	"github.com/Sanpeta/rate-limiter-pos-go-expert/internal/limiter"
)

// RateLimitMiddleware cria um middleware para limitar a taxa de requisições.
func RateLimitMiddleware(limiter limiter.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// Verificar se o limite de taxa foi excedido
			allowed, retryAfter := limiter.Allow(ctx, r)

			// Se excedido, HTTP 429 Too Many Requests
			if !allowed {
				w.Header().Set("Retry-After", fmt.Sprintf("%.0f", retryAfter.Seconds()))
				http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}

			// Continua para próxima
			next.ServeHTTP(w, r)
		})
	}
}
