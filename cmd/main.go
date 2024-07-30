package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Sanpeta/rate-limiter-pos-go-expert/internal/config"
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

}
