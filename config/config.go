package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Server   ServerConfig
	Caixa    CaixaConfig
	Service  ServiceConfig
	Database DatabaseConfig
}

type AppConfig struct {
	Name        string
	Version     string
	Environment string
}

type ServerConfig struct {
	Port            string
	ReadTimeout     int
	WriteTimeout    int
	MaxConnections  int
	GracefulTimeout int
}

type CaixaConfig struct {
	BaseURL            string
	Timeout            int
	RateLimitPerSec    int
	CacheTTLMinutes    int
	MaxRetries         int
	RetryBackoffMs     int
	UserAgent          string
}

type ServiceConfig struct {
	DownstreamURL    string
	DownstreamTimeout int
	AsyncForwarding  bool
	MaxQueueSize     int
}

type DatabaseConfig struct {
	Enabled    bool
	Type       string
	DSN        string
	MaxConns   int
	IdleConns  int
	ConnMaxAge int
}

func Load() (*Config, error) {
	// Carregar .env se existir
	_ = godotenv.Load(".env", ".env.local")

	return &Config{
		App: AppConfig{
			Name:        getEnv("APP_NAME", "lottery-caixa-service"),
			Version:     getEnv("APP_VERSION", "1.0.0"),
			Environment: getEnv("ENVIRONMENT", "development"),
		},
		Server: ServerConfig{
			Port:            getEnv("PORT", "8080"),
			ReadTimeout:     getEnvInt("READ_TIMEOUT", 15),
			WriteTimeout:    getEnvInt("WRITE_TIMEOUT", 15),
			MaxConnections:  getEnvInt("MAX_CONNECTIONS", 1000),
			GracefulTimeout: getEnvInt("GRACEFUL_TIMEOUT", 30),
		},
		Caixa: CaixaConfig{
			BaseURL:         getEnv("CAIXA_BASE_URL", "http://servicebus2.caixa.gov.br/portaldeloterias/api/"),
			Timeout:         getEnvInt("CAIXA_TIMEOUT", 15),
			RateLimitPerSec: getEnvInt("CAIXA_RATE_LIMIT", 2),
			CacheTTLMinutes: getEnvInt("CAIXA_CACHE_TTL", 5),
			MaxRetries:      getEnvInt("CAIXA_MAX_RETRIES", 3),
			RetryBackoffMs:  getEnvInt("CAIXA_RETRY_BACKOFF", 1000),
			UserAgent:       getEnv("CAIXA_USER_AGENT", "LotteryService/1.0"),
		},
		Service: ServiceConfig{
			DownstreamURL:    getEnv("DOWNSTREAM_SERVICE_URL", ""),
			DownstreamTimeout: getEnvInt("DOWNSTREAM_TIMEOUT", 10),
			AsyncForwarding:  getEnvBool("ASYNC_FORWARDING", true),
			MaxQueueSize:     getEnvInt("MAX_QUEUE_SIZE", 1000),
		},
		Database: DatabaseConfig{
			Enabled:    getEnvBool("DATABASE_ENABLED", false),
			Type:       getEnv("DATABASE_TYPE", "mysql"),
			DSN:        getEnv("DATABASE_DSN", ""),
			MaxConns:   getEnvInt("DATABASE_MAX_CONNS", 25),
			IdleConns:  getEnvInt("DATABASE_IDLE_CONNS", 5),
			ConnMaxAge: getEnvInt("DATABASE_CONN_MAX_AGE", 5),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return value == "true" || value == "1" || value == "yes"
	}
	return defaultValue
}

func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("PORT não pode estar vazio")
	}

	if c.Caixa.BaseURL == "" {
		return fmt.Errorf("CAIXA_BASE_URL não pode estar vazio")
	}

	if c.Service.DownstreamURL == "" && c.Service.AsyncForwarding {
		// Aviso: forwarding async sem URL configurada
		fmt.Println("WARNING: DOWNSTREAM_SERVICE_URL não configurada, async forwarding desativado")
	}

	return nil
}

func (c *Config) String() string {
	return fmt.Sprintf(
		"App: %s v%s\nEnvironment: %s\nServer: :%s\nCaixa URL: %s\nDownstream: %s",
		c.App.Name, c.App.Version, c.App.Environment, c.Server.Port, c.Caixa.BaseURL, c.Service.DownstreamURL,
	)
}
