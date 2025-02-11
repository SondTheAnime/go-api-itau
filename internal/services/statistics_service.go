package services

import (
	"math"
	"sync"
	"time"

	"api-itau/config"
	"api-itau/handlers"
	"api-itau/internal/models"
)

type StatisticsService struct {
	transactions []models.Transaction
	windowSize   time.Duration
	mu           sync.RWMutex
}

func NewStatisticsService(cfg *config.Config) *StatisticsService {
	return &StatisticsService{
		transactions: make([]models.Transaction, 0),
		windowSize:   time.Duration(cfg.Stats.WindowSeconds) * time.Second,
	}
}

func (s *StatisticsService) AddTransaction(t models.Transaction) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cleanOldTransactions()
	s.transactions = append(s.transactions, t)
}

func (s *StatisticsService) GetStatistics() (*handlers.StatisticsResponse, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.cleanOldTransactions()

	if len(s.transactions) == 0 {
		return nil, nil
	}

	var sum float64
	min := math.MaxFloat64
	max := -math.MaxFloat64

	for _, t := range s.transactions {
		sum += t.Value
		if t.Value < min {
			min = t.Value
		}
		if t.Value > max {
			max = t.Value
		}
	}

	count := len(s.transactions)
	avg := sum / float64(count)

	return &handlers.StatisticsResponse{
		Count: count,
		Sum:   sum,
		Avg:   avg,
		Min:   min,
		Max:   max,
	}, nil
}

func (s *StatisticsService) DeleteTransactions() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.transactions = make([]models.Transaction, 0)
}

func (s *StatisticsService) cleanOldTransactions() {
	now := time.Now()
	cutoff := now.Add(-s.windowSize)

	validTransactions := make([]models.Transaction, 0)
	for _, t := range s.transactions {
		if t.Timestamp.After(cutoff) {
			validTransactions = append(validTransactions, t)
		}
	}

	s.transactions = validTransactions
}
