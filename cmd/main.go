package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Sanpeta/rate-limiter-pos-go-expert/internal/config"
	"github.com/Sanpeta/rate-limiter-pos-go-expert/internal/limiter"
	"github.com/Sanpeta/rate-limiter-pos-go-expert/internal/limiter/strategy"
	"github.com/Sanpeta/rate-limiter-pos-go-expert/internal/middleware"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Carregar configurações
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Error loading configuration from .: %v", err)
	}

	fmt.Print(config)

	// Criar cliente Redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.RedisAddr,
	})

	// Testar a conexão com o Redis
	err = redisClient.Ping(context.Background()).Err()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis successfully!")

	// Criar IPLimiter com a estratégia Redis
	ipLimiter := limiter.NewIPLimiter(config, strategy.NewRedisStrategy(redisClient))

	// Criar TokenLimiter com a estratégia Redis
	tokenLimiter := limiter.NewTokenLimiter(config, strategy.NewRedisStrategy(redisClient))

	// Escolher o limitador a ser usado
	var limiterToUse limiter.Limiter
	if config.LimitBy == "ip" {
		limiterToUse = ipLimiter
	} else {
		limiterToUse = tokenLimiter
	}

	// Criar o roteador e adicionar o middleware
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Aplicar o middleware ao roteador
	handlerWithMiddleware := middleware.RateLimitMiddleware(limiterToUse)(mux)

	// Iniciar o servidor
	log.Fatal(http.ListenAndServe(":8080", handlerWithMiddleware))

}
