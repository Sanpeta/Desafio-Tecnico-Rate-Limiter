package main

import (
	"log"

	"github.com/Sanpeta/rate-limiter-pos-go-expert/internal/config"
)

func main() {
	// Carregar configurações
	_, err := config.LoadConfig("../")
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

}
