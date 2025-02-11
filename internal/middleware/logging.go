package middleware

import (
	"api-itau/pkg/logger"
	"context"
	"net/http"
	"time"
)

// responseWriter é um wrapper para http.ResponseWriter que captura o status code
type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.status = code
	rw.wroteHeader = true
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

// LoggingMiddleware cria um middleware para logging de requisições HTTP
func LoggingMiddleware(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Captura informações da requisição
			path := r.URL.Path
			method := r.Method
			remoteAddr := r.RemoteAddr
			userAgent := r.UserAgent()

			// Log da requisição recebida
			log.Info("requisição recebida",
				"method", method,
				"path", path,
				"remote_addr", remoteAddr,
				"user_agent", userAgent,
			)

			// Wrapper para capturar o status code
			rw := newResponseWriter(w)

			// Processa a requisição
			next.ServeHTTP(rw, r)

			// Calcula a duração
			duration := time.Since(start)

			// Log da resposta
			log.Info("requisição processada",
				"method", method,
				"path", path,
				"status", rw.Status(),
				"duration", duration.String(),
				"remote_addr", remoteAddr,
			)
		})
	}
}

// RecoveryMiddleware cria um middleware para recuperação de pânico
func RecoveryMiddleware(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					log.Error("pânico recuperado",
						"error", err,
						"path", r.URL.Path,
						"method", r.Method,
					)
					http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// RequestIDMiddleware adiciona um ID único para cada requisição
func RequestIDMiddleware(log logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("X-Request-ID")
			if requestID == "" {
				requestID = generateRequestID()
			}

			// Adiciona o ID no header da resposta
			w.Header().Set("X-Request-ID", requestID)

			// Adiciona o ID no contexto
			ctx := r.Context()
			ctx = setRequestID(ctx, requestID)
			r = r.WithContext(ctx)

			log.Info("request-id atribuído",
				"request_id", requestID,
				"path", r.URL.Path,
			)

			next.ServeHTTP(w, r)
		})
	}
}

// generateRequestID gera um ID único para a requisição
func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(6)
}

// randomString gera uma string aleatória de tamanho n
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}

// Chave do contexto para o request ID
type contextKey string

const requestIDKey = contextKey("requestID")

// setRequestID adiciona o request ID no contexto
func setRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// GetRequestID obtém o request ID do contexto
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(requestIDKey).(string); ok {
		return requestID
	}
	return ""
}
