package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// LotteryPayload estrutura esperada do serviço upstream
type LotteryPayload struct {
	Status   string `json:"status"`
	Results  []struct {
		ID          string    `json:"id"`
		GameType    string    `json:"gameType"`
		DrawNumber  int       `json:"drawNumber"`
		DrawDate    string    `json:"drawDate"`
		Numbers     []int     `json:"numbers"`
		Winners     int       `json:"winners"`
		Prize       float64   `json:"prize"`
		ProcessedAt time.Time `json:"processedAt"`
	} `json:"results"`
	Metadata struct {
		ProcessedAt   time.Time `json:"processedAt"`
		SourceService string    `json:"sourceService"`
		TotalRecords  int       `json:"totalRecords"`
		RequestID     string    `json:"requestId"`
	} `json:"metadata"`
	ErrorMsg string `json:"errorMsg,omitempty"`
}

// HandleProcessLottery processa resultados de loteria
func HandleProcessLottery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload LotteryPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Printf("Erro ao decodificar payload: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	requestID := payload.Metadata.RequestID
	log.Printf("[%s] Recebido payload com %d resultados", requestID, payload.Metadata.TotalRecords)

	// Validar payload
	if payload.Status == "error" {
		log.Printf("[%s] Erro recebido: %s", requestID, payload.ErrorMsg)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "error",
			"error":  payload.ErrorMsg,
		})
		return
	}

	// Processar cada resultado
	for _, result := range payload.Results {
		log.Printf("[%s] Processando %s - Sorteio %d",
			requestID, result.GameType, result.DrawNumber)

		// Aqui você implementaria sua lógica:
		// - Validar números
		// - Armazenar em banco de dados
		// - Calcular estatísticas
		// - Enviar notificações
		// etc.

		if len(result.Numbers) > 0 {
			log.Printf("[%s] Números: %v", requestID, result.Numbers)
		}
	}

	// Responder sucesso
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   "processed",
		"records":  payload.Metadata.TotalRecords,
		"requestId": requestID,
		"timestamp": time.Now().UTC(),
	})

	log.Printf("[%s] Processamento concluído", requestID)
}

// HandleHealth retorna status de saúde
func HandleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
		"service": "lottery-processor",
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	http.HandleFunc("/api/lottery/process", HandleProcessLottery)
	http.HandleFunc("/health", HandleHealth)

	log.Printf("Serviço de processamento iniciado na porta %s", port)
	log.Printf("Endpoints:")
	log.Printf("  POST http://localhost:%s/api/lottery/process", port)
	log.Printf("  GET  http://localhost:%s/health", port)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
