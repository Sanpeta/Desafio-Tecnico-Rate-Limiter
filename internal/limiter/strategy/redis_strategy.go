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

// Allow verifica se uma requisição é permitida e atualiza o contador no Redis.
func (redis *RedisStrategy) Allow(ctx context.Context, key string, maxRequests int, expire time.Duration) (bool, time.Duration) {
	// Gerar a chave no Redis com base no IP ou Token (por exemplo, "ratelimit:192.168.1.1" ou "ratelimit:abc123")
	timestamp := time.Now().Unix()
	redisKey := fmt.Sprintf("ratelimit:%s:%d", key, timestamp)

	log.Printf("Before INCR: key=%s, expire=%s", redisKey, expire)

	// Incrementar o contador no Redis (com expiração)
	count, err := redis.client.Incr(ctx, redisKey).Result()
	if err != nil {
		log.Printf("Error incrementing key: %v", err)
		return false, 0
	}

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
			log.Printf("Error getting TTL: %v", err)
			return false, 0
		}
		return false, ttl
	}

	// Retornar true (permitido) e 0 (sem tempo de bloqueio)
	return true, 0
}
