package ratelimit

import (
	"fmt"
	"sync"
	"time"
)

// RateLimiter implementa rate limiting com token bucket
type RateLimiter struct {
	ticker     *time.Ticker
	ch         chan struct{}
	maxBurst   int
	rate       int
	mu         sync.Mutex
	lastReset  time.Time
	requests   int
}

// NewRateLimiter cria um novo rate limiter
// requestsPerSecond: número de requisições permitidas por segundo
func NewRateLimiter(requestsPerSecond int) *RateLimiter {
	if requestsPerSecond <= 0 {
		requestsPerSecond = 1
	}

	rl := &RateLimiter{
		rate:      requestsPerSecond,
		maxBurst:  requestsPerSecond * 2,
		ch:        make(chan struct{}, requestsPerSecond*2),
		lastReset: time.Now(),
	}

	// Preencher com tokens iniciais
	for i := 0; i < requestsPerSecond; i++ {
		rl.ch <- struct{}{}
	}

	// Refill tokens periodicamente
	go rl.refill()

	return rl
}

// Wait aguarda até que uma requisição possa ser feita
func (rl *RateLimiter) Wait() {
	<-rl.ch
}

// TryWait tenta fazer uma requisição sem bloquear
func (rl *RateLimiter) TryWait() bool {
	select {
	case <-rl.ch:
		return true
	default:
		return false
	}
}

// WaitWithTimeout aguarda com timeout
func (rl *RateLimiter) WaitWithTimeout(timeout time.Duration) bool {
	select {
	case <-rl.ch:
		return true
	case <-time.After(timeout):
		return false
	}
}

// refill recarrega tokens periodicamente
func (rl *RateLimiter) refill() {
	ticker := time.NewTicker(time.Second / time.Duration(rl.rate))
	defer ticker.Stop()

	for range ticker.C {
		select {
		case rl.ch <- struct{}{}:
		default:
			// Buffer cheio, ignorar
		}
	}
}

// GetStats retorna estatísticas do rate limiter
func (rl *RateLimiter) GetStats() map[string]interface{} {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	return map[string]interface{}{
		"rate":           rl.rate,
		"maxBurst":       rl.maxBurst,
		"availableTokens": len(rl.ch),
		"requests":       rl.requests,
	}
}

// Close encerra o rate limiter
func (rl *RateLimiter) Close() {
	// Goroutine de refill será finalizada automaticamente
}

// WindowedRateLimiter implementa rate limiting por janela deslizante
type WindowedRateLimiter struct {
	maxRequests int
	window      time.Duration
	requests    []time.Time
	mu          sync.Mutex
}

// NewWindowedRateLimiter cria um novo rate limiter de janela
func NewWindowedRateLimiter(maxRequests int, window time.Duration) *WindowedRateLimiter {
	return &WindowedRateLimiter{
		maxRequests: maxRequests,
		window:      window,
		requests:    make([]time.Time, 0, maxRequests),
	}
}

// Allow verifica se uma requisição é permitida
func (wrl *WindowedRateLimiter) Allow() bool {
	wrl.mu.Lock()
	defer wrl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-wrl.window)

	// Remover requisições antigas
	validRequests := make([]time.Time, 0)
	for _, req := range wrl.requests {
		if req.After(cutoff) {
			validRequests = append(validRequests, req)
		}
	}
	wrl.requests = validRequests

	// Verificar limite
	if len(wrl.requests) < wrl.maxRequests {
		wrl.requests = append(wrl.requests, now)
		return true
	}

	return false
}

// AllowWithError retorna erro se não permitido
func (wrl *WindowedRateLimiter) AllowWithError() error {
	if !wrl.Allow() {
		return fmt.Errorf("rate limit exceeded: máximo %d requisições por %v", 
			wrl.maxRequests, wrl.window)
	}
	return nil
}

// GetStatus retorna status do rate limiter
func (wrl *WindowedRateLimiter) GetStatus() map[string]interface{} {
	wrl.mu.Lock()
	defer wrl.mu.Unlock()

	return map[string]interface{}{
		"maxRequests": wrl.maxRequests,
		"window":      wrl.window.String(),
		"currentCount": len(wrl.requests),
		"remaining":   wrl.maxRequests - len(wrl.requests),
	}
}

// Reset limpa as requisições
func (wrl *WindowedRateLimiter) Reset() {
	wrl.mu.Lock()
	defer wrl.mu.Unlock()
	wrl.requests = make([]time.Time, 0, wrl.maxRequests)
}
