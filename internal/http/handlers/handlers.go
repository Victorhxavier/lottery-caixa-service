package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/victor/lottery-caixa-service/internal/domain"
	"github.com/victor/lottery-caixa-service/internal/service"
)

func StartVerifyScheduler() {
	ticker := time.NewTicker(2 * time.Hour)
	go func() {
		for range ticker.C {
			callVerifyEndpoint()
		}
	}()
	callVerifyEndpoint()
}

func callVerifyEndpoint() {
	resp, err := http.Get("https://backend-loteria-xpn7.onrender.com/api/verify")
	if err != nil {
		log.Printf("Erro ao chamar /api/verify: %v", err)
		return
	}
	defer resp.Body.Close()
	log.Printf("Chamada /api/verify status: %d", resp.StatusCode)
}

// GetLotteryResults retorna resultados de loteria diretamente da API da Caixa
func GetLotteryResults(svc *service.LotteryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		gameType := r.URL.Query().Get("gameType")
		if gameType == "" {
			gameType = "lotofacil"
		}

		concurso := r.URL.Query().Get("concurso")

		requestID := uuid.New().String()
		log.Printf("[%s] Buscando resultados para: %s (concurso: %s)", requestID, gameType, concurso)

		result, err := svc.FetchLotteryResults(r.Context(), gameType, concurso)
		if err != nil {
			log.Printf("[%s] Erro ao buscar resultados: %v", requestID, err)
			respondError(w, http.StatusInternalServerError, "Erro ao buscar resultados", err, requestID)
			return
		}

		respondJSON(w, http.StatusOK, result)
	}
}

// WebhookResults processa webhook de resultados
func WebhookResults(svc *service.LotteryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			respondError(w, http.StatusMethodNotAllowed, "Método não permitido", nil, "")
			return
		}

		requestID := uuid.New().String()
		var webhook domain.WebhookPayload

		if err := json.NewDecoder(r.Body).Decode(&webhook); err != nil {
			log.Printf("[%s] Erro ao decodificar webhook: %v", requestID, err)
			respondError(w, http.StatusBadRequest, "Payload inválido", err, requestID)
			return
		}

		result := &domain.LotteryResult{
			ID:          uuid.New().String(),
			GameType:    webhook.GameType,
			DrawNumber:  webhook.DrawNumber,
			DrawDate:    webhook.DrawDate,
			Numbers:     webhook.Numbers,
			ProcessedAt: time.Now().UTC(),
			Source:      "webhook",
		}

		payload := svc.CreateDownstreamPayload([]domain.LotteryResult{*result})
		payload.Metadata.RequestID = requestID

		if err := svc.SendToDownstream(r.Context(), payload); err != nil {
			log.Printf("[%s] Erro ao processar webhook: %v", requestID, err)
			respondError(w, http.StatusInternalServerError, "Erro ao processar", err, requestID)
			return
		}

		respondJSON(w, http.StatusAccepted, map[string]interface{}{
			"status":    "processed",
			"requestId": requestID,
		})
	}
}

// HealthCheck verifica saúde do serviço
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := domain.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	respondJSON(w, http.StatusOK, response)
}

// ReadinessCheck verifica se o serviço está pronto
func ReadinessCheck(svc *service.LotteryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !svc.IsHealthy() {
			respondJSON(w, http.StatusServiceUnavailable, map[string]string{
				"status": "not ready",
			})
			return
		}

		respondJSON(w, http.StatusOK, map[string]string{
			"status": "ready",
		})
	}
}

// Metrics retorna métricas do serviço
func Metrics(w http.ResponseWriter, r *http.Request) {
	// Placeholder - será implementado em próxima versão
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("# HELP lottery_service_up Service uptime indicator\n"))
	w.Write([]byte("# TYPE lottery_service_up gauge\n"))
	w.Write([]byte("lottery_service_up 1\n"))
}

// ServiceInfo retorna informações do serviço
func ServiceInfo(w http.ResponseWriter, r *http.Request) {
	// Seria preenchido com informações do serviço
	respondJSON(w, http.StatusOK, map[string]string{
		"name":    "lottery-caixa-service",
		"version": "1.0.0",
		"status":  "ok",
	})
}

// respondJSON escreve uma resposta JSON
func respondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// respondError escreve uma resposta de erro
func respondError(w http.ResponseWriter, statusCode int, message string, err error, requestID string) {
	errMsg := message
	if err != nil {
		errMsg = message + ": " + err.Error()
	}

	response := domain.ErrorResponse{
		Error:      message,
		Message:    errMsg,
		StatusCode: statusCode,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
		RequestID:  requestID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}
