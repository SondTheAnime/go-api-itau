package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// Logger define a interface para logging
type Logger interface {
	Info(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})
}

// DefaultLogger é a implementação padrão do Logger
type DefaultLogger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

// NewDefaultLogger cria uma nova instância do DefaultLogger
func NewDefaultLogger() *DefaultLogger {
	return &DefaultLogger{
		infoLogger:  log.New(os.Stdout, "[INFO] ", log.LstdFlags),
		errorLogger: log.New(os.Stderr, "[ERROR] ", log.LstdFlags),
	}
}

// formatKeyvals formata os pares chave-valor em uma string
func formatKeyvals(keyvals ...interface{}) string {
	if len(keyvals)%2 != 0 {
		keyvals = append(keyvals, "MISSING")
	}

	var pairs []string
	for i := 0; i < len(keyvals); i += 2 {
		key := fmt.Sprintf("%v", keyvals[i])
		value := fmt.Sprintf("%v", keyvals[i+1])
		pairs = append(pairs, fmt.Sprintf("%s=%s", key, value))
	}

	return strings.Join(pairs, " ")
}

// Info registra uma mensagem de informação
func (l *DefaultLogger) Info(msg string, keyvals ...interface{}) {
	l.infoLogger.Printf("%s %s", msg, formatKeyvals(keyvals...))
}

// Error registra uma mensagem de erro
func (l *DefaultLogger) Error(msg string, keyvals ...interface{}) {
	l.errorLogger.Printf("%s %s", msg, formatKeyvals(keyvals...))
}
