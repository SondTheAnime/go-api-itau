package tests

import (
	"bytes"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"api-itau/config"
	"api-itau/handlers"
	"api-itau/internal/models"
	"api-itau/internal/services"
	"api-itau/pkg/utils"
)

// mockLogger é uma implementação mock do logger para testes
type mockLogger struct{}

func (l *mockLogger) Info(msg string, keyvals ...interface{})  {}
func (l *mockLogger) Error(msg string, keyvals ...interface{}) {}

// setupTimeProvider configura um provedor de tempo mockado para testes
func setupTimeProvider() (*utils.MockTimeProvider, *config.Config) {
	cfg := &config.Config{
		Stats: config.StatsConfig{WindowSeconds: 60},
	}
	mockTime := utils.NewMockTimeProvider(time.Now())
	utils.SetTimeProvider(mockTime)
	return mockTime, cfg
}

// floatEquals verifica se dois números float64 são aproximadamente iguais
func floatEquals(a, b float64) bool {
	const epsilon = 0.01
	return math.Abs(a-b) < epsilon
}

// TestTransactionEndpoints testa os endpoints de transação
func TestTransactionEndpoints(t *testing.T) {
	mockTime, cfg := setupTimeProvider()
	log := &mockLogger{}

	statsService := services.NewStatisticsService(cfg, log)
	transactionService := services.NewTransactionService(statsService, log)
	handler := handlers.NewTransactionHandler(transactionService, log)

	baseTime := mockTime.Now()

	tests := []struct {
		name           string
		method         string
		body           map[string]interface{}
		expectedStatus int
	}{
		{
			name:   "Transação válida",
			method: http.MethodPost,
			body: map[string]interface{}{
				"valor":    123.45,
				"dataHora": baseTime.Add(-30 * time.Second).Format(time.RFC3339),
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:   "Transação com valor negativo",
			method: http.MethodPost,
			body: map[string]interface{}{
				"valor":    -10.00,
				"dataHora": baseTime.Format(time.RFC3339),
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:   "Transação no futuro",
			method: http.MethodPost,
			body: map[string]interface{}{
				"valor":    100.00,
				"dataHora": baseTime.Add(1 * time.Hour).Format(time.RFC3339),
			},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name:           "Deletar transações",
			method:         http.MethodDelete,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request

			if tt.method == http.MethodPost {
				bodyBytes, _ := json.Marshal(tt.body)
				req = httptest.NewRequest(tt.method, "/transacao", bytes.NewBuffer(bodyBytes))
			} else {
				req = httptest.NewRequest(tt.method, "/transacao", nil)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler retornou status code errado: obtido %v esperado %v",
					status, tt.expectedStatus)
			}
		})
	}
}

// TestStatisticsEndpoint testa o endpoint de estatísticas
func TestStatisticsEndpoint(t *testing.T) {
	mockTime, cfg := setupTimeProvider()
	log := &mockLogger{}

	statsService := services.NewStatisticsService(cfg, log)
	handler := handlers.NewStatisticsHandler(statsService, log)

	baseTime := mockTime.Now()

	// Adiciona algumas transações
	transactions := []models.Transaction{
		{Value: 100.00, Timestamp: baseTime.Add(-30 * time.Second)},
		{Value: 50.00, Timestamp: baseTime.Add(-45 * time.Second)},
		{Value: 25.00, Timestamp: baseTime.Add(-15 * time.Second)},
	}

	for _, tx := range transactions {
		statsService.AddTransaction(tx)
	}

	// Testa o endpoint GET /estatistica
	t.Run("Obter estatísticas", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/estatistica", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler retornou status code errado: obtido %v esperado %v",
				status, http.StatusOK)
		}

		var response handlers.StatisticsResponse
		err := json.NewDecoder(rr.Body).Decode(&response)
		if err != nil {
			t.Fatalf("erro ao decodificar resposta: %v", err)
		}

		// Verifica os valores esperados
		expectedStats := handlers.StatisticsResponse{
			Count: 3,
			Sum:   175.00,
			Avg:   58.33,
			Min:   25.00,
			Max:   100.00,
		}

		if response.Count != expectedStats.Count {
			t.Errorf("count incorreto: obtido %v esperado %v",
				response.Count, expectedStats.Count)
		}

		if response.Sum != expectedStats.Sum {
			t.Errorf("sum incorreto: obtido %v esperado %v",
				response.Sum, expectedStats.Sum)
		}

		if !floatEquals(response.Avg, expectedStats.Avg) {
			t.Errorf("avg incorreto: obtido %v esperado %v",
				response.Avg, expectedStats.Avg)
		}

		if response.Min != expectedStats.Min {
			t.Errorf("min incorreto: obtido %v esperado %v",
				response.Min, expectedStats.Min)
		}

		if response.Max != expectedStats.Max {
			t.Errorf("max incorreto: obtido %v esperado %v",
				response.Max, expectedStats.Max)
		}
	})

	// Testa estatísticas após janela de tempo
	t.Run("Estatísticas após expiração", func(t *testing.T) {
		// Avança o tempo em 61 segundos
		novoTempo := mockTime.Now().Add(61 * time.Second)
		mockTime.Set(novoTempo)

		// Limpa as transações antigas
		statsService.DeleteTransactions()

		// Adiciona uma nova transação com o tempo atual
		novaTransacao := models.Transaction{
			Value:     1.00,
			Timestamp: mockTime.Now(),
		}
		statsService.AddTransaction(novaTransacao)

		req := httptest.NewRequest(http.MethodGet, "/estatistica", nil)
		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		var response handlers.StatisticsResponse
		json.NewDecoder(rr.Body).Decode(&response)

		// Deve ter apenas a nova transação
		if response.Count != 1 {
			t.Errorf("deveria ter 1 transação, mas tem %d", response.Count)
		}

		// Verifica se os valores estão corretos
		if response.Sum != 1.00 {
			t.Errorf("soma incorreta: obtido %v esperado %v", response.Sum, 1.00)
		}
	})
}

// TestConcurrency testa o comportamento da API em cenários concorrentes
func TestConcurrency(t *testing.T) {
	mockTime, cfg := setupTimeProvider()
	log := &mockLogger{}

	statsService := services.NewStatisticsService(cfg, log)
	transactionService := services.NewTransactionService(statsService, log)
	transactionHandler := handlers.NewTransactionHandler(transactionService, log)
	statisticsHandler := handlers.NewStatisticsHandler(statsService, log)

	baseTime := mockTime.Now()

	// Cria vários requests concorrentes
	const numRequests = 100
	done := make(chan bool)

	for i := 0; i < numRequests; i++ {
		go func(i int) {
			// Alterna entre POST e GET
			if i%2 == 0 {
				body := map[string]interface{}{
					"valor":    float64(i) + 0.99,
					"dataHora": baseTime.Add(-time.Duration(i) * time.Second).Format(time.RFC3339),
				}
				bodyBytes, _ := json.Marshal(body)
				req := httptest.NewRequest(http.MethodPost, "/transacao", bytes.NewBuffer(bodyBytes))
				rr := httptest.NewRecorder()
				transactionHandler.ServeHTTP(rr, req)
			} else {
				req := httptest.NewRequest(http.MethodGet, "/estatistica", nil)
				rr := httptest.NewRecorder()
				statisticsHandler.ServeHTTP(rr, req)
			}
			done <- true
		}(i)
	}

	// Aguarda todas as goroutines terminarem
	for i := 0; i < numRequests; i++ {
		<-done
	}

	// Verifica se as estatísticas estão consistentes
	req := httptest.NewRequest(http.MethodGet, "/estatistica", nil)
	rr := httptest.NewRecorder()
	statisticsHandler.ServeHTTP(rr, req)

	var response handlers.StatisticsResponse
	json.NewDecoder(rr.Body).Decode(&response)

	if response.Count != numRequests/2 {
		t.Errorf("número incorreto de transações: obtido %v esperado %v",
			response.Count, numRequests/2)
	}
}
