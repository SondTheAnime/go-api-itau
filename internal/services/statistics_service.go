package services

import (
	"math"
	"sync"
	"time"

	"api-itau/config"
	"api-itau/handlers"
	"api-itau/internal/models"
	"api-itau/pkg/logger"
	"api-itau/pkg/utils"
)

// StatisticsService implementa a interface handlers.StatisticsService
type StatisticsService struct {
	transactions []models.Transaction
	window       *utils.SlidingWindow
	mu           sync.RWMutex
	logger       logger.Logger
}

// NewStatisticsService cria uma nova instância do StatisticsService
func NewStatisticsService(cfg *config.Config, log logger.Logger) *StatisticsService {
	duration := time.Duration(cfg.Stats.WindowSeconds) * time.Second
	window := utils.NewSlidingWindow(duration, utils.GetTimeProvider())

	return &StatisticsService{
		transactions: make([]models.Transaction, 0),
		window:       window,
		logger:       log,
	}
}

// AddTransaction adiciona uma nova transação
func (s *StatisticsService) AddTransaction(t models.Transaction) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cleanOldTransactions()
	s.transactions = append(s.transactions, t)

	s.logger.Info("transação adicionada às estatísticas",
		"valor", t.Value,
		"dataHora", t.Timestamp,
	)
}

// GetStatistics retorna as estatísticas das transações dentro da janela de tempo
func (s *StatisticsService) GetStatistics() (*handlers.StatisticsResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cleanOldTransactions()

	if len(s.transactions) == 0 {
		return &handlers.StatisticsResponse{
			Count: 0,
			Sum:   0,
			Avg:   0,
			Min:   0,
			Max:   0,
		}, nil
	}

	stats := s.calculateStatistics(s.transactions)

	s.logger.Info("estatísticas calculadas",
		"count", stats.Count,
		"sum", stats.Sum,
		"avg", stats.Avg,
		"min", stats.Min,
		"max", stats.Max,
	)

	return stats, nil
}

// DeleteTransactions remove todas as transações
func (s *StatisticsService) DeleteTransactions() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.transactions = make([]models.Transaction, 0)
	s.logger.Info("todas as transações foram removidas das estatísticas")
}

// cleanOldTransactions remove transações fora da janela de tempo
func (s *StatisticsService) cleanOldTransactions() {
	validTransactions := make([]models.Transaction, 0)
	for _, t := range s.transactions {
		if s.window.IsInWindow(t.Timestamp) {
			validTransactions = append(validTransactions, t)
		}
	}

	if len(s.transactions) != len(validTransactions) {
		s.logger.Info("transações antigas removidas",
			"removidas", len(s.transactions)-len(validTransactions),
		)
	}

	s.transactions = validTransactions
}

// calculateStatistics calcula as estatísticas para um conjunto de transações
func (s *StatisticsService) calculateStatistics(transactions []models.Transaction) *handlers.StatisticsResponse {
	var sum float64
	min := math.MaxFloat64
	max := -math.MaxFloat64

	for _, t := range transactions {
		sum += t.Value
		if t.Value < min {
			min = t.Value
		}
		if t.Value > max {
			max = t.Value
		}
	}

	count := len(transactions)
	avg := sum / float64(count)

	return &handlers.StatisticsResponse{
		Count: count,
		Sum:   sum,
		Avg:   avg,
		Min:   min,
		Max:   max,
	}
}
