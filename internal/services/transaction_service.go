package services

import (
	"api-itau/internal/models"
	"api-itau/pkg/logger"
)

// TransactionService implementa a interface handlers.TransactionService
type TransactionService struct {
	statsService *StatisticsService
	logger       logger.Logger
}

// NewTransactionService cria uma nova instância do TransactionService
func NewTransactionService(statsService *StatisticsService, logger logger.Logger) *TransactionService {
	return &TransactionService{
		statsService: statsService,
		logger:       logger,
	}
}

// AddTransaction adiciona uma nova transação
func (s *TransactionService) AddTransaction(t models.Transaction) error {
	// Adiciona a transação ao serviço de estatísticas
	s.statsService.AddTransaction(t)

	s.logger.Info("transação adicionada com sucesso",
		"valor", t.Value,
		"dataHora", t.Timestamp,
	)

	return nil
}

// DeleteTransactions remove todas as transações
func (s *TransactionService) DeleteTransactions() error {
	// Remove as transações do serviço de estatísticas
	s.statsService.DeleteTransactions()

	s.logger.Info("todas as transações foram removidas")

	return nil
}
