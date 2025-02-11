package handlers

import (
	"encoding/json"
	"net/http"

	"api-itau/pkg/logger"
)

type StatisticsResponse struct {
	Count int     `json:"count"`
	Sum   float64 `json:"sum"`
	Avg   float64 `json:"avg"`
	Min   float64 `json:"min"`
	Max   float64 `json:"max"`
}

type StatisticsService interface {
	GetStatistics() (*StatisticsResponse, error)
}

// StatisticsHandler encapsula a lógica de manipulação de requisições de estatísticas
type StatisticsHandler struct {
	service StatisticsService
	logger  logger.Logger
}

// NewStatisticsHandler cria uma nova instância do StatisticsHandler
func NewStatisticsHandler(service StatisticsService, logger logger.Logger) *StatisticsHandler {
	return &StatisticsHandler{
		service: service,
		logger:  logger,
	}
}

// ServeHTTP implementa a interface http.Handler
func (h *StatisticsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Verifica se o método é GET
	if r.Method != http.MethodGet {
		h.logger.Error("método não permitido", "método", r.Method)
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	stats, err := h.service.GetStatistics()
	if err != nil {
		h.logger.Error("erro ao obter estatísticas", "erro", err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}

	if stats == nil {
		stats = &StatisticsResponse{
			Count: 0,
			Sum:   0,
			Avg:   0,
			Min:   0,
			Max:   0,
		}
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(stats); err != nil {
		h.logger.Error("erro ao codificar resposta", "erro", err)
		http.Error(w, "Erro ao processar resposta", http.StatusInternalServerError)
		return
	}

	h.logger.Info("estatísticas retornadas com sucesso",
		"count", stats.Count,
		"sum", stats.Sum,
		"avg", stats.Avg,
		"min", stats.Min,
		"max", stats.Max,
	)
}

type Logger interface {
	Info(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})
}
