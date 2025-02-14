package handlers

import (
	"net/http"

	"api-itau/pkg/logger"
)

// StatisticsResponse representa a resposta com as estatísticas das transações
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
		RespondWithError(w, http.StatusMethodNotAllowed, "method_not_allowed", "Método não permitido")
		return
	}

	stats, err := h.service.GetStatistics()
	if err != nil {
		h.logger.Error("erro ao obter estatísticas", "erro", err)
		RespondWithError(w, http.StatusInternalServerError, "internal_error", "Erro interno do servidor")
		return
	}

	// Garante que temos um objeto de estatísticas válido
	if stats == nil {
		stats = &StatisticsResponse{
			Count: 0,
			Sum:   0,
			Avg:   0,
			Min:   0,
			Max:   0,
		}
	}

	h.logger.Info("estatísticas retornadas com sucesso",
		"count", stats.Count,
		"sum", stats.Sum,
		"avg", stats.Avg,
		"min", stats.Min,
		"max", stats.Max,
	)

	// Cria uma resposta formatada
	response := APIResponse{
		Success: true,
		Data:    stats,
	}

	// Envia a resposta
	RespondWithJSON(w, http.StatusOK, response)
}

type Logger interface {
	Info(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})
}
