package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"api-itau/internal/models"
	"api-itau/pkg/logger"
)

// TransactionRequest representa o payload da requisição de transação
type TransactionRequest struct {
	Value     float64   `json:"valor"`
	Timestamp time.Time `json:"dataHora"`
}

// TransactionService define o contrato para o serviço de transações
type TransactionService interface {
	AddTransaction(models.Transaction) error
	DeleteTransactions() error
}

// TransactionHandler encapsula a lógica de manipulação de requisições de transações
type TransactionHandler struct {
	service TransactionService
	logger  logger.Logger
}

// NewTransactionHandler cria uma nova instância do TransactionHandler
func NewTransactionHandler(service TransactionService, logger logger.Logger) *TransactionHandler {
	return &TransactionHandler{
		service: service,
		logger:  logger,
	}
}

// ServeHTTP implementa a interface http.Handler
func (h *TransactionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.handlePost(w, r)
	case http.MethodDelete:
		h.handleDelete(w, r)
	default:
		h.logger.Error("método não permitido", "método", r.Method)
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

// handlePost processa requisições POST para criar uma nova transação
func (h *TransactionHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	// Limita o tamanho do corpo da requisição para prevenir ataques
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20)) // 1 MB
	if err != nil {
		h.logger.Error("erro ao ler corpo da requisição", "erro", err)
		http.Error(w, "Erro ao ler requisição", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req TransactionRequest
	if err := json.Unmarshal(body, &req); err != nil {
		h.logger.Error("erro ao decodificar JSON", "erro", err)
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Cria e valida a transação
	transaction, err := models.NewTransaction(req.Value, req.Timestamp)
	if err != nil {
		h.logger.Error("transação inválida", "erro", err)
		http.Error(w, "Transação inválida", http.StatusUnprocessableEntity)
		return
	}

	// Adiciona a transação através do serviço
	if err := h.service.AddTransaction(*transaction); err != nil {
		h.logger.Error("erro ao adicionar transação", "erro", err)
		http.Error(w, "Erro ao processar transação", http.StatusInternalServerError)
		return
	}

	h.logger.Info("transação criada com sucesso",
		"valor", transaction.Value,
		"dataHora", transaction.Timestamp,
	)

	w.WriteHeader(http.StatusCreated)
}

// handleDelete processa requisições DELETE para remover todas as transações
func (h *TransactionHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	if err := h.service.DeleteTransactions(); err != nil {
		h.logger.Error("erro ao deletar transações", "erro", err)
		http.Error(w, "Erro ao deletar transações", http.StatusInternalServerError)
		return
	}

	h.logger.Info("todas as transações foram deletadas")
	w.WriteHeader(http.StatusOK)
}
