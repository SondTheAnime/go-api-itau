package models

import (
	"fmt"
	"time"
)

type Transaction struct {
	Value     float64   `json:"valor"`
	Timestamp time.Time `json:"dataHora"`
}

func (t *Transaction) Validate() error {
	if t.Value < 0 {
		return fmt.Errorf("Valor não pode ser negativo")
	}

	if t.Timestamp.After(time.Now()) {
		return	fmt.Errorf("data da transação não pode estar no futuro")
	}

	return nil
}

func NewTransaction(value float64, timestamp time.Time) (*Transaction, error) {
	t := &Transaction{
		Value:     value,
		Timestamp: timestamp,
	}

	if err := t.Validate(); err != nil {
		return nil, fmt.Errorf("erro ao criar transação: %w", err)
	}

	return t, nil
}