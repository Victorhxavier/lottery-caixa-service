package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/victor/lottery-caixa-service/config"
	"github.com/victor/lottery-caixa-service/internal/cache"
	"github.com/victor/lottery-caixa-service/internal/domain"
	"github.com/victor/lottery-caixa-service/internal/ratelimit"
)

// LotteryService gerencia operações de loteria
type LotteryService struct {
	cfg           *config.Config
	httpClient    *http.Client
	rateLimiter   *ratelimit.RateLimiter
	cache         *cache.MemoryCache
	forwardQueue  chan *domain.DownstreamPayload
	metrics       *ServiceMetrics
	startTime     time.Time
	isHealthy     atomic.Bool
}

// ServiceMetrics contém métricas do serviço
type ServiceMetrics struct {
	requestsTotal   atomic.Int64
	requestsSuccess atomic.Int64
	requestsError   atomic.Int64
	cacheHits       atomic.Int64
	cacheMisses     atomic.Int64
	totalLatency    atomic.Int64
	downstreamSent  atomic.Int64
	mu              sync.RWMutex
	latencies       []time.Duration
}

// NewLotteryService cria um novo serviço de loteria
func NewLotteryService(cfg *config.Config) *LotteryService {
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Erro ao validar configuração: %v", err)
	}

	svc := &LotteryService{
		cfg:       cfg,
		cache:     cache.NewMemoryCache(cfg.Caixa.CacheTTLMinutes),
		rateLimiter: ratelimit.NewRateLimiter(cfg.Caixa.RateLimitPerSec),
		metrics:   &ServiceMetrics{},
		startTime: time.Now(),
		httpClient: &http.Client{
			Timeout: time.Duration(cfg.Caixa.Timeout) * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
				DisableCompression:  false,
			},
		},
	}

	svc.isHealthy.Store(true)

	// Inicializar fila de forwarding
	if cfg.Service.AsyncForwarding && cfg.Service.DownstreamURL != "" {
		svc.forwardQueue = make(chan *domain.DownstreamPayload, cfg.Service.MaxQueueSize)
		go svc.processForwardingQueue()
	}

	log.Printf("LotteryService iniciado: %s", cfg)

	return svc
}

// FetchLotteryResults busca resultados da loteria com retry e retorna objeto da Caixa
func (s *LotteryService) FetchLotteryResults(ctx context.Context, gameType string, concurso string) (*domain.CaixaAPIResponse, error) {
	// Incrementar métrica
	s.metrics.requestsTotal.Add(1)

	// Verificar cache
	cacheKey := gameType
	if concurso != "" {
		cacheKey = fmt.Sprintf("%s:%s", gameType, concurso)
	}
	if cached, exists := s.cache.Get(cacheKey); exists {
		s.metrics.cacheHits.Add(1)
		log.Printf("Cache hit para %s", cacheKey)
		if result, ok := cached.(*domain.CaixaAPIResponse); ok {
			return result, nil
		}
	}

	s.metrics.cacheMisses.Add(1)

	// Implementar retry com backoff exponencial
	var lastErr error
	for attempt := 0; attempt < s.cfg.Caixa.MaxRetries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(math.Pow(2, float64(attempt))) *
				time.Duration(s.cfg.Caixa.RetryBackoffMs) * time.Millisecond
			log.Printf("Retry %d após %v para %s", attempt+1, backoff, cacheKey)
			time.Sleep(backoff)
		}

		result, err := s.fetchWithRateLimit(ctx, gameType, concurso)
		if err == nil {
			s.metrics.requestsSuccess.Add(1)
			s.cache.Set(cacheKey, result)
			return result, nil
		}

		lastErr = err
	}

	s.metrics.requestsError.Add(1)
	s.isHealthy.Store(false)
	return nil, fmt.Errorf("falha ao buscar resultados após %d tentativas: %w",
		s.cfg.Caixa.MaxRetries, lastErr)
}

// fetchWithRateLimit busca resultados respeitando rate limit
func (s *LotteryService) fetchWithRateLimit(ctx context.Context, gameType string, concurso string) (*domain.CaixaAPIResponse, error) {
	start := time.Now()

	// Aguardar rate limit
	s.rateLimiter.Wait()

	defer func() {
		latency := time.Since(start)
		s.metrics.totalLatency.Add(int64(latency))
		s.metrics.mu.Lock()
		s.metrics.latencies = append(s.metrics.latencies, latency)
		s.metrics.mu.Unlock()
	}()

	// Construir URL - se concurso for informado, adiciona à URL
	var url string
	if concurso != "" {
		url = fmt.Sprintf("%s%s/%s", s.cfg.Caixa.BaseURL, gameType, concurso)
	} else {
		// API da Caixa sem número busca o último concurso
		url = fmt.Sprintf("%s%s", s.cfg.Caixa.BaseURL, gameType)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("User-Agent", s.cfg.Caixa.UserAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("erro ao chamar API da Caixa: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API retornou status %d: %s", resp.StatusCode, string(body))
	}

	var apiResp domain.CaixaAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta: %w", err)
	}

	return &apiResp, nil
}

// SendToDownstream envia resultados para o serviço downstream
func (s *LotteryService) SendToDownstream(ctx context.Context, payload *domain.DownstreamPayload) error {
	if s.cfg.Service.DownstreamURL == "" {
		return fmt.Errorf("DOWNSTREAM_SERVICE_URL não configurada")
	}

	if s.cfg.Service.AsyncForwarding {
		select {
		case s.forwardQueue <- payload:
			return nil
		default:
			return fmt.Errorf("fila de forwarding cheia")
		}
	}

	return s.sendToDownstreamSync(ctx, payload)
}

// sendToDownstreamSync envia sincronamente
func (s *LotteryService) sendToDownstreamSync(ctx context.Context, payload *domain.DownstreamPayload) error {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("erro ao serializar payload: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 
		time.Duration(s.cfg.Service.DownstreamTimeout)*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", s.cfg.Service.DownstreamURL, 
		bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("erro ao criar requisição: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Service-Source", s.cfg.App.Name)
	req.Header.Set("X-Request-ID", payload.Metadata.RequestID)

	client := &http.Client{Timeout: time.Duration(s.cfg.Service.DownstreamTimeout) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao enviar para downstream: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("downstream retornou status %d: %s", resp.StatusCode, string(body))
	}

	s.metrics.downstreamSent.Add(1)
	log.Printf("Resultado enviado para downstream: %s", payload.Metadata.RequestID)
	return nil
}

// processForwardingQueue processa a fila de forwarding async
func (s *LotteryService) processForwardingQueue() {
	for payload := range s.forwardQueue {
		if err := s.sendToDownstreamSync(context.Background(), payload); err != nil {
			log.Printf("Erro ao enviar para downstream: %v", err)
		}
	}
}

// GetServiceInfo retorna informações do serviço
func (s *LotteryService) GetServiceInfo() *domain.ServiceInfo {
	return &domain.ServiceInfo{
		Name:        s.cfg.App.Name,
		Version:     s.cfg.App.Version,
		Environment: s.cfg.App.Environment,
		Uptime:      time.Since(s.startTime).String(),
		Status:      map[bool]string{true: "healthy", false: "unhealthy"}[s.isHealthy.Load()],
	}
}

// IsHealthy verifica se o serviço está saudável
func (s *LotteryService) IsHealthy() bool {
	return s.isHealthy.Load()
}

// GetMetrics retorna métricas do serviço
func (s *LotteryService) GetMetrics() *domain.MetricsResponse {
	s.metrics.mu.RLock()
	defer s.metrics.mu.RUnlock()

	var avgLatency time.Duration
	if len(s.metrics.latencies) > 0 {
		total := int64(0)
		for _, l := range s.metrics.latencies {
			total += int64(l)
		}
		avgLatency = time.Duration(total / int64(len(s.metrics.latencies)))
	}

	return &domain.MetricsResponse{
		RequestsTotal:  s.metrics.requestsTotal.Load(),
		RequestsSuccess: s.metrics.requestsSuccess.Load(),
		RequestsError:  s.metrics.requestsError.Load(),
		AverageLatency: avgLatency,
		CacheHits:      s.metrics.cacheHits.Load(),
		CacheMisses:    s.metrics.cacheMisses.Load(),
	}
}

// CreateDownstreamPayload cria um payload para envio downstream
func (s *LotteryService) CreateDownstreamPayload(results []domain.LotteryResult) *domain.DownstreamPayload {
	return &domain.DownstreamPayload{
		Status:  "success",
		Results: results,
		Metadata: domain.Metadata{
			ProcessedAt:   time.Now().UTC(),
			SourceService: s.cfg.App.Name,
			TotalRecords:  len(results),
			RequestID:     uuid.New().String(),
		},
	}
}

// Close fecha recursos do serviço
func (s *LotteryService) Close() {
	log.Printf("Encerrando LotteryService...")
	
	if s.forwardQueue != nil {
		close(s.forwardQueue)
	}
	
	s.rateLimiter.Close()
	log.Printf("LotteryService encerrado")
}
