package utils

import (
	"sync"
	"time"
)

var (
	defaultProvider TimeProvider = &RealTimeProvider{}
	currentProvider TimeProvider = defaultProvider
	providerMu      sync.RWMutex
)

// SetTimeProvider define o provedor de tempo global
func SetTimeProvider(provider TimeProvider) {
	providerMu.Lock()
	defer providerMu.Unlock()
	currentProvider = provider
}

// GetTimeProvider retorna o provedor de tempo atual
func GetTimeProvider() TimeProvider {
	providerMu.RLock()
	defer providerMu.RUnlock()
	return currentProvider
}

// ResetTimeProvider restaura o provedor de tempo padrão
func ResetTimeProvider() {
	SetTimeProvider(defaultProvider)
}

// TimeProvider é uma interface para obter o tempo atual
// Útil para testes e para garantir consistência temporal
type TimeProvider interface {
	Now() time.Time
}

// RealTimeProvider é a implementação padrão que usa o tempo real do sistema
type RealTimeProvider struct{}

func (p *RealTimeProvider) Now() time.Time {
	return time.Now()
}

// MockTimeProvider é uma implementação para testes que permite controlar o tempo
type MockTimeProvider struct {
	mu      sync.RWMutex
	current time.Time
}

func NewMockTimeProvider(t time.Time) *MockTimeProvider {
	return &MockTimeProvider{current: t}
}

func (p *MockTimeProvider) Now() time.Time {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.current
}

func (p *MockTimeProvider) Set(t time.Time) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.current = t
}

func (p *MockTimeProvider) Add(d time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.current = p.current.Add(d)
}

// TimeWindow representa uma janela de tempo com início e fim
type TimeWindow struct {
	Start time.Time
	End   time.Time
}

// NewTimeWindow cria uma nova janela de tempo
func NewTimeWindow(start, end time.Time) TimeWindow {
	return TimeWindow{
		Start: start,
		End:   end,
	}
}

// Contains verifica se um timestamp está dentro da janela de tempo
func (w TimeWindow) Contains(t time.Time) bool {
	return (t.Equal(w.Start) || t.After(w.Start)) && (t.Equal(w.End) || t.Before(w.End))
}

// Duration retorna a duração da janela de tempo
func (w TimeWindow) Duration() time.Duration {
	return w.End.Sub(w.Start)
}

// SlidingWindow representa uma janela de tempo deslizante
type SlidingWindow struct {
	duration time.Duration
	provider TimeProvider
}

// NewSlidingWindow cria uma nova janela deslizante
func NewSlidingWindow(duration time.Duration, provider TimeProvider) *SlidingWindow {
	if provider == nil {
		provider = GetTimeProvider()
	}
	return &SlidingWindow{
		duration: duration,
		provider: provider,
	}
}

// GetWindow retorna a janela de tempo atual
func (w *SlidingWindow) GetWindow() TimeWindow {
	now := w.provider.Now()
	return TimeWindow{
		Start: now.Add(-w.duration),
		End:   now,
	}
}

// IsInWindow verifica se um timestamp está dentro da janela atual
func (w *SlidingWindow) IsInWindow(t time.Time) bool {
	return w.GetWindow().Contains(t)
}

// FormatISO formata um time.Time no padrão ISO 8601
func FormatISO(t time.Time) string {
	return t.Format(time.RFC3339)
}

// ParseISO faz o parse de uma string ISO 8601 para time.Time
func ParseISO(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}
