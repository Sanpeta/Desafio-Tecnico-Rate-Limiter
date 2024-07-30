package limiter

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/Sanpeta/rate-limiter-pos-go-expert/internal/config"
	"github.com/Sanpeta/rate-limiter-pos-go-expert/internal/limiter/strategy"
)

// IPLimiter implementa a interface Limiter para limitar requisições por IP.
type IPLimiter struct {
	storage strategy.RedisStrategy
	config  config.Config
}

// NewIPLimiter cria uma nova instância de IPLimiter.
func NewIPLimiter(config config.Config, storage *strategy.RedisStrategy) *IPLimiter {
	return &IPLimiter{
		storage: *storage,
		config:  config,
	}
}

// Retorna true se permitida, false se bloqueada, e a duração do bloqueio de um determinado IP.
func (ipLimiter *IPLimiter) Allow(ctx context.Context, r *http.Request) (bool, time.Duration) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return false, 0
	}

	// Obter o limite máximo de requisições e a duração do bloqueio da configuração
	maxRequests := ipLimiter.config.MaxRequestsPerSecondIP
	blockDuration := ipLimiter.config.BlockDuration

	// Usar a estratégia de armazenamento para verificar e atualizar o limite
	allowed, remainingTime := ipLimiter.storage.Allow(ctx, ip, maxRequests, blockDuration)

	// retornar false e a duração restante do bloqueio
	if !allowed {
		return false, remainingTime
	}

	// retornar true e 0 (sem bloqueio)
	return true, 0
}
