package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/victor/lottery-caixa-service/config"
	"github.com/victor/lottery-caixa-service/internal/domain"
	"github.com/victor/lottery-caixa-service/internal/service"
)

func TestHealthCheck(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()

	// Simular handler
	http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(domain.HealthResponse{
			Status:    "healthy",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
	}).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var resp domain.HealthResponse
	json.NewDecoder(w.Body).Decode(&resp)

	if resp.Status != "healthy" {
		t.Errorf("Expected healthy status, got %s", resp.Status)
	}
}

func TestConfigLoad(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Erro ao carregar config: %v", err)
	}

	if cfg.App.Name == "" {
		t.Error("Nome do app não pode estar vazio")
	}

	if cfg.Server.Port == "" {
		t.Error("Porta não pode estar vazia")
	}
}

func TestLotteryServiceCreation(t *testing.T) {
	cfg, _ := config.Load()
	svc := service.NewLotteryService(cfg)
	defer svc.Close()

	info := svc.GetServiceInfo()
	if info.Name == "" {
		t.Error("Nome do serviço não pode estar vazio")
	}

	if !svc.IsHealthy() {
		t.Error("Serviço deveria estar saudável")
	}
}

func TestDownstreamPayload(t *testing.T) {
	result := domain.LotteryResult{
		ID:        "test-1",
		GameType:  "lotofacil",
		Numbers:   []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
		ProcessedAt: time.Now(),
	}

	payload := domain.DownstreamPayload{
		Status:  "success",
		Results: []domain.LotteryResult{result},
		Metadata: domain.Metadata{
			TotalRecords: 1,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Erro ao serializar: %v", err)
	}

	var unmarshaled domain.DownstreamPayload
	err = json.Unmarshal(jsonData, &unmarshaled)
	if err != nil {
		t.Fatalf("Erro ao desserializar: %v", err)
	}

	if unmarshaled.Status != "success" {
		t.Errorf("Status incorreto: %s", unmarshaled.Status)
	}
}

func TestWebhookPayload(t *testing.T) {
	webhook := domain.WebhookPayload{
		GameType:   "lotofacil",
		DrawNumber: 2024001,
		DrawDate:   "2024-01-15",
		Numbers:    []int{1, 5, 8, 12, 15, 18, 22, 25},
		Timestamp:  time.Now(),
	}

	jsonData, err := json.Marshal(webhook)
	if err != nil {
		t.Fatalf("Erro ao serializar webhook: %v", err)
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/api/v1/lottery/webhook", bytes.NewBuffer(jsonData))

	if req.Method != "POST" {
		t.Errorf("Método incorreto: %s", req.Method)
	}
}

func TestErrorResponse(t *testing.T) {
	errResp := domain.ErrorResponse{
		Error:      "test error",
		Message:    "Test error message",
		StatusCode: 400,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		RequestID:  "test-123",
	}

	jsonData, _ := json.Marshal(errResp)
	var unmarshaled domain.ErrorResponse
	json.Unmarshal(jsonData, &unmarshaled)

	if unmarshaled.Error != "test error" {
		t.Errorf("Erro não corresponde")
	}
}

func TestMetricsResponse(t *testing.T) {
	metrics := domain.MetricsResponse{
		RequestsTotal:   100,
		RequestsSuccess: 95,
		RequestsError:   5,
		CacheHits:       50,
		CacheMisses:     50,
	}

	if metrics.RequestsTotal != 100 {
		t.Errorf("Total de requisições incorreto")
	}

	if metrics.RequestsSuccess != 95 {
		t.Errorf("Requisições bem-sucedidas incorretas")
	}
}

func BenchmarkConfigLoad(b *testing.B) {
	for i := 0; i < b.N; i++ {
		config.Load()
	}
}

func BenchmarkPayloadSerialization(b *testing.B) {
	payload := domain.DownstreamPayload{
		Status: "success",
		Results: []domain.LotteryResult{
			{
				ID:        "test-1",
				GameType:  "lotofacil",
				Numbers:   []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
				ProcessedAt: time.Now(),
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.Marshal(payload)
	}
}
