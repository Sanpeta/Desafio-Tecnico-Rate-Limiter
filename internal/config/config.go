package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config é para armazena as configurações do rate limiter.
type Config struct {
	LimitBy                   string        `mapstructure:"LIMIT_BY"`
	MaxRequestsPerSecondIP    int           `mapstructure:"MAX_REQUESTS_PER_SECOND_IP"`
	MaxRequestsPerSecondToken int           `mapstructure:"MAX_REQUESTS_PER_SECOND_TOKEN"`
	BlockDuration             time.Duration `mapstructure:"BLOCK_DURATION_SECONDS"`
	RedisAddr                 string        `mapstructure:"REDIS_ADDR"`
}

// LoadConfig carrega as configurações do arquivo .env e variáveis de ambiente.
func LoadConfig(path string) (config Config, err error) {
	// Configurar o Viper
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	// Ler o arquivo de configuração
	err = viper.ReadInConfig()
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}

	// Mapear as configurações para a struct
	err = viper.Unmarshal(&config)
	if err != nil {
		return config, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}
