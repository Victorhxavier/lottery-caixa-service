package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/victor/lottery-caixa-service/config"
	"github.com/victor/lottery-caixa-service/internal/http/handlers"
	"github.com/victor/lottery-caixa-service/internal/http/middleware"
	"github.com/victor/lottery-caixa-service/internal/service"
)

func main() {
	// Carregar configurações
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}

	log.Printf("Iniciando Lottery Caixa Service na porta %s", cfg.Server.Port)
	log.Printf("Modo: %s", cfg.App.Environment)

	// Inicializar serviço
	lotteryService := service.NewLotteryService(cfg)
	defer lotteryService.Close()

	// Criar roteador
	mux := setupRouter(lotteryService)

	// Criar servidor HTTP
	srv := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Canal para erros
	errChan := make(chan error, 1)

	// Iniciar servidor em goroutine
	go func() {
		log.Printf("Servidor HTTP iniciado em %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- fmt.Errorf("erro ao iniciar servidor: %w", err)
		}
	}()

	// Aguardar sinais de shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errChan:
		log.Printf("Erro: %v", err)
	case sig := <-sigChan:
		log.Printf("Sinal recebido: %v", sig)

		// Graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		log.Printf("Encerrando servidor...")
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Erro ao encerrar servidor: %v", err)
		}
	}

	log.Printf("Servidor encerrado")
}

func setupRouter(svc *service.LotteryService) *http.ServeMux {
	mux := http.NewServeMux()

	// Aplicar middlewares globais
	handler := middleware.Chain(
		middleware.Logging,
		middleware.Recovery,
		middleware.CORS,
	)

	// Rotas de saúde
	mux.HandleFunc("/health", handlers.HealthCheck)
	mux.HandleFunc("/ready", handlers.ReadinessCheck(svc))

	// Rotas de métricas
	mux.HandleFunc("/metrics", handlers.Metrics)

	// Rotas de negócio
	mux.HandleFunc("/", handler(handlers.GetLotteryResults(svc)))
	mux.HandleFunc("/api/v1/lottery", handler(handlers.GetLotteryResults(svc)))
	mux.HandleFunc("/api/v1/lottery/webhook", handler(handlers.WebhookResults(svc)))

	// Rotas de informações
	mux.HandleFunc("/api/v1/info", handler(handlers.ServiceInfo))

	return mux
}
