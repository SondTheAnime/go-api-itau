package validator

import (
	"fmt"
	"time"
)

// TransactionValidator encapsula a lógica de validação de transações
type TransactionValidator struct {
	maxValue float64
}

// NewTransactionValidator cria uma nova instância do validador
func NewTransactionValidator(maxValue float64) *TransactionValidator {
	return &TransactionValidator{
		maxValue: maxValue,
	}
}

// ValidateValue verifica se o valor da transação é válido
func (v *TransactionValidator) ValidateValue(value float64) error {
	if value < 0 {
		return fmt.Errorf("valor não pode ser negativo")
	}

	if value > v.maxValue {
		return fmt.Errorf("valor excede o limite máximo permitido de %.2f", v.maxValue)
	}

	return nil
}

// ValidateTimestamp verifica se o timestamp da transação é válido
func (v *TransactionValidator) ValidateTimestamp(timestamp time.Time) error {
	now := time.Now()

	// Verifica se a data está no futuro
	if timestamp.After(now) {
		return fmt.Errorf("data da transação não pode estar no futuro")
	}

	// Verifica se a data é muito antiga (mais de 5 anos)
	if timestamp.Before(now.AddDate(-5, 0, 0)) {
		return fmt.Errorf("data da transação é muito antiga")
	}

	return nil
}

// ValidateJSON verifica se os campos obrigatórios estão presentes
func (v *TransactionValidator) ValidateJSON(hasValue, hasTimestamp bool) error {
	if !hasValue {
		return fmt.Errorf("campo 'valor' é obrigatório")
	}

	if !hasTimestamp {
		return fmt.Errorf("campo 'dataHora' é obrigatório")
	}

	return nil
}

// IsValidISOTimestamp verifica se a string está no formato ISO 8601
func IsValidISOTimestamp(timestamp string) bool {
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05.999Z07:00",
	}

	for _, layout := range layouts {
		if _, err := time.Parse(layout, timestamp); err == nil {
			return true
		}
	}

	return false
}

// ParseTimestamp tenta fazer o parse do timestamp em vários formatos
func ParseTimestamp(timestamp string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05.999Z07:00",
	}

	var firstErr error
	for _, layout := range layouts {
		t, err := time.Parse(layout, timestamp)
		if err == nil {
			return t, nil
		}
		if firstErr == nil {
			firstErr = err
		}
	}

	return time.Time{}, fmt.Errorf("formato de data inválido: %w", firstErr)
}
