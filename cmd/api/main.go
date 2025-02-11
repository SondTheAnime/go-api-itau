package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultPort               = "8080"
	defaultStatsWindowSeconds = 60
)

func main() {

	logger := log.New(os.Stdout, "[API] ", log.LstdFlags|log.Lshortfile)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy"}`))
	})

	mux.HandleFunc("POST /transacao", nil)
	mux.HandleFunc("DELETE /transacao", nil)
	mux.HandleFunc("GET /estatistica", nil)

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	serverErrors := make(chan error, 1)

	go func() {
		logger.Printf("Servidor está escutando na porta: %s...", port)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		logger.Fatalf("erro ao iniciar o servidor: %v", err)

	case sig := <-shutdown:
		logger.Printf("iniciando shutdown, sinal: %v", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Tenta desligar o servidor graciosamente
		if err := server.Shutdown(ctx); err != nil {
			logger.Printf("erro durante o shutdown do servidor: %v", err)
			if err := server.Close(); err != nil {
				logger.Printf("erro ao forçar fechamento do servidor: %v", err)
			}
		}

		logger.Println("servidor desligado com sucesso")
	}
}
