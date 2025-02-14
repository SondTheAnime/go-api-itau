package handlers

import (
	"encoding/json"
	"net/http"
)

// APIResponse é a estrutura base para todas as respostas da API
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError representa um erro na API
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// RespondWithJSON envia uma resposta JSON com o status code apropriado
func RespondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{
			Success: false,
			Error: &APIError{
				Code:    "internal_error",
				Message: "Erro ao processar resposta",
			},
		})
	}
}

// RespondWithError envia uma resposta de erro padronizada
func RespondWithError(w http.ResponseWriter, statusCode int, code string, message string) {
	RespondWithJSON(w, statusCode, APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
		},
	})
}

// RespondWithSuccess envia uma resposta de sucesso padronizada
func RespondWithSuccess(w http.ResponseWriter, statusCode int, data interface{}) {
	RespondWithJSON(w, statusCode, APIResponse{
		Success: true,
		Data:    data,
	})
}
