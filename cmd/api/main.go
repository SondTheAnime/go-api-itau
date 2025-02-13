package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"api-itau/config"
	"api-itau/handlers"
	"api-itau/internal/middleware"
	"api-itau/internal/services"
	"api-itau/pkg/logger"

	scalar "github.com/MarceloPetrucio/go-scalar-api-reference"
)

func main() {
	// Carrega as configurações
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

	// Inicializa o logger
	log := logger.NewDefaultLogger()

	// Cria os serviços
	statsService := services.NewStatisticsService(cfg, log)
	transactionService := services.NewTransactionService(statsService, log)

	// Cria os handlers
	statsHandler := handlers.NewStatisticsHandler(statsService, log)
	transactionHandler := handlers.NewTransactionHandler(transactionService, log)

	// Cria o router
	mux := http.NewServeMux()

	// Registra as rotas
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	mux.Handle("POST /transacao", transactionHandler)
	mux.Handle("DELETE /transacao", transactionHandler)
	mux.Handle("GET /estatistica", statsHandler)

	// Adiciona a rota para a documentação
	mux.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
		htmlContent, err := scalar.ApiReferenceHTML(&scalar.Options{
			SpecURL: "./docs/scalar/scalar.yaml",
			CustomOptions: scalar.CustomOptions{
				PageTitle: "API de Transações - Documentação",
			},
			DarkMode: true,
		})
		if err != nil {
			log.Error("erro ao gerar documentação", "erro", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(htmlContent))
	})

	// Aplica os middlewares
	handler := middleware.RequestIDMiddleware(log)(
		middleware.LoggingMiddleware(log)(
			middleware.RecoveryMiddleware(log)(mux),
		),
	)

	// Cria o servidor
	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Canal para erros do servidor
	serverErrors := make(chan error, 1)

	// Inicia o servidor em uma goroutine
	go func() {
		log.Info("servidor iniciado", "porta", cfg.Server.Port)
		serverErrors <- server.ListenAndServe()
	}()

	// Canal para sinais de interrupção
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Aguarda sinal de shutdown ou erro
	select {
	case err := <-serverErrors:
		log.Error("erro ao iniciar servidor", "erro", err)

	case sig := <-shutdown:
		log.Info("iniciando shutdown", "sinal", sig)

		// Contexto com timeout para shutdown gracioso
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Tenta desligar o servidor graciosamente
		if err := server.Shutdown(ctx); err != nil {
			log.Error("erro durante shutdown do servidor", "erro", err)
			if err := server.Close(); err != nil {
				log.Error("erro ao forçar fechamento do servidor", "erro", err)
			}
		}

		log.Info("servidor desligado com sucesso")
	}
}
