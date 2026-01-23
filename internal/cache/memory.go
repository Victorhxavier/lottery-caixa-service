package cache

import (
	"sync"
	"time"

	"github.com/victor/lottery-caixa-service/internal/domain"
)

// CacheEntry representa um item no cache
type CacheEntry struct {
	Data      interface{}
	ExpiresAt time.Time
}

// MemoryCache implementa cache em memória com TTL
type MemoryCache struct {
	data  map[string]*CacheEntry
	mu    sync.RWMutex
	ttl   time.Duration
	hits  int64
	misses int64
}

// NewMemoryCache cria um novo cache em memória
func NewMemoryCache(ttlMinutes int) *MemoryCache {
	mc := &MemoryCache{
		data: make(map[string]*CacheEntry),
		ttl:  time.Duration(ttlMinutes) * time.Minute,
	}

	// Limpar expirados periodicamente
	go mc.cleanupExpired()

	return mc
}

// Get retorna um valor do cache
func (mc *MemoryCache) Get(key string) (interface{}, bool) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	entry, exists := mc.data[key]
	if !exists {
		mc.misses++
		return nil, false
	}

	if time.Now().After(entry.ExpiresAt) {
		mc.misses++
		return nil, false
	}

	mc.hits++
	return entry.Data, true
}

// Set armazena um valor no cache
func (mc *MemoryCache) Set(key string, value interface{}) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.data[key] = &CacheEntry{
		Data:      value,
		ExpiresAt: time.Now().Add(mc.ttl),
	}
}

// Delete remove um valor do cache
func (mc *MemoryCache) Delete(key string) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	delete(mc.data, key)
}

// Clear limpa todo o cache
func (mc *MemoryCache) Clear() {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	mc.data = make(map[string]*CacheEntry)
	mc.hits = 0
	mc.misses = 0
}

// GetLottery retorna um resultado de loteria do cache
func (mc *MemoryCache) GetLottery(gameType string) (*domain.LotteryResult, bool) {
	data, exists := mc.Get(gameType)
	if !exists {
		return nil, false
	}

	if result, ok := data.(*domain.LotteryResult); ok {
		return result, true
	}

	return nil, false
}

// SetLottery armazena um resultado de loteria no cache
func (mc *MemoryCache) SetLottery(gameType string, result *domain.LotteryResult) {
	mc.Set(gameType, result)
}

// GetStats retorna estatísticas do cache
func (mc *MemoryCache) GetStats() map[string]interface{} {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	total := mc.hits + mc.misses
	var hitRate float64
	if total > 0 {
		hitRate = float64(mc.hits) / float64(total)
	}

	return map[string]interface{}{
		"size":      len(mc.data),
		"hits":      mc.hits,
		"misses":    mc.misses,
		"hitRate":   hitRate,
		"ttl":       mc.ttl.String(),
	}
}

// cleanupExpired remove entradas expiradas periodicamente
func (mc *MemoryCache) cleanupExpired() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		mc.mu.Lock()

		now := time.Now()
		for key, entry := range mc.data {
			if now.After(entry.ExpiresAt) {
				delete(mc.data, key)
			}
		}

		mc.mu.Unlock()
	}
}

// Size retorna o número de itens no cache
func (mc *MemoryCache) Size() int {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	return len(mc.data)
}
