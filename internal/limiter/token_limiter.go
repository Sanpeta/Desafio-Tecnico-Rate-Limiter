package limiter

import (
	"context"
	"net/http"
	"time"

	"github.com/Sanpeta/rate-limiter-pos-go-expert/internal/config"
	"github.com/Sanpeta/rate-limiter-pos-go-expert/internal/limiter/strategy"
)

// TokenLimiter implementa a interface Limiter para limitar requisições por token.
type TokenLimiter struct {
	storage strategy.RedisStrategy
	config  config.Config
}

// NewTokenLimiter cria uma nova instância de TokenLimiter.
func NewTokenLimiter(config config.Config, storage *strategy.RedisStrategy) *TokenLimiter {
	return &TokenLimiter{
		storage: *storage,
		config:  config,
	}
}

// Retorna true se permitida, false se bloqueada, e a duração do bloqueio de um determinado TOKEN.
func (l *TokenLimiter) Allow(ctx context.Context, r *http.Request) (bool, time.Duration) {
	// Extrair o token do cabeçalho da requisição
	token := r.Header.Get("API_KEY")
	if token == "" {
		return false, 0
	}

	// Obter o limite máximo de requisições e a duração do bloqueio da configuração
	maxRequests := l.config.MaxRequestsPerSecondToken
	blockDuration := l.config.BlockDuration

	// Usar a estratégia de armazenamento para verificar e atualizar o limite
	allowed, remainingTime := l.storage.Allow(ctx, token, maxRequests, blockDuration)

	// Se não permitido, retornar false e a duração restante do bloqueio
	if !allowed {
		return false, remainingTime
	}

	// Se permitido, retornar true e 0 (sem bloqueio)
	return true, 0
}
