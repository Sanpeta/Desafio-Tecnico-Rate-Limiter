package strategy

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisStrategy implementa a interface Strategy usando Redis como armazenamento.
type RedisStrategy struct {
	client *redis.Client
}

// NewRedisStrategy cria uma nova instância de RedisStrategy.
func NewRedisStrategy(client *redis.Client) *RedisStrategy {
	return &RedisStrategy{client: client}
}

// Allow erifica se uma requisição é permitida e atualiza o contador no Redis.
func (redis *RedisStrategy) Allow(ctx context.Context, ip string, maxRequests int, expire time.Duration) (bool, time.Duration) {
	// Gerar a chave no Redis com base no IP e na janela de tempo (por exemplo, "ratelimit:192.168.1.1:2024072920")
	redisKey := fmt.Sprintf("ratelimit:%s:%d", ip, time.Now().Unix()/int64(expire.Seconds()))

	// Log antes do incremento
	log.Printf("Before INCR: key=%s, expire=%s", redisKey, expire)

	// Incrementar o contador no Redis (com expiração)
	count, err := redis.client.Incr(ctx, redisKey).Result()
	if err != nil {
		return false, 0
	}

	// Log após o incremento
	log.Printf("After INCR: key=%s, count=%d, expire=%s", redisKey, count, expire)

	// Definir a expiração da chave se for a primeira requisição na janela de tempo
	if count == 1 {
		redis.client.Expire(ctx, redisKey, expire)
	}

	// Verificar se o limite foi excedido
	if int(count) > maxRequests {
		// Obter o tempo restante para a expiração da chave
		ttl, err := redis.client.TTL(ctx, redisKey).Result()
		if err != nil {
			return false, 0
		}
		return false, ttl
	}

	// Retornar true (permitido) e 0 (sem tempo de bloqueio)
	return true, 0
}
