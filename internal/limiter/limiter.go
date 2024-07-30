package limiter

import (
	"context"
	"net/http"
	"time"
)

// Limiter é a interface que define os métodos necessários para um rate limiter.
type Limiter interface {
	Allow(ctx context.Context, r *http.Request) (bool, time.Duration)
}
