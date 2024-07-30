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
	// Chave para contagem de requisições
	redisKey := fmt.Sprintf("ratelimit:%s", key)
	// Chave para bloqueio
	blockedKey := fmt.Sprintf("ratelimit:%s:block", key)

	// Verificar se o IP/token está bloqueado
	blocked, err := redis.client.Exists(ctx, blockedKey).Result()
	if err != nil {
		log.Println("Redis error on EXISTS: ", err)
		return false, 0
	}
	if blocked == 1 {
		ttl, err := redis.client.TTL(ctx, blockedKey).Result()
		if err != nil {
			log.Println("Redis error on TTL: ", err)
			return false, 0
		}
		return false, ttl
	}

	// Incrementar o contador no Redis (com expiração)
	count, err := redis.client.Incr(ctx, redisKey).Result()
	if err != nil {
		log.Println("Error incrementing key: ", err)
		return false, 0
	}

	// Definido a expiração da chave se for a primeira requisição na janela de tempo
	if count == 1 {
		err = redis.client.Expire(ctx, redisKey, time.Second).Err()
		if err != nil {
			log.Println("Error setting expiration: ", err)
			return false, 0
		}
	}

	// Verifica se o limite foi excedido
	if int(count) > maxRequests {
		_, err := redis.client.Set(ctx, blockedKey, 1, expire).Result()
		// Bloquear por 5 minutos se o limite for excedido
		if err != nil {
			log.Println("Redis error on SET: ", err)
		}
		return false, expire
	}

	// Permitido
	return true, 0
}
