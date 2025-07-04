package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type ProcessStats struct {
	Total     int `json:"total"`
	Active    int `json:"active"`
	Paused    int `json:"paused"`
	Archived  int `json:"archived"`
	ThisMonth int `json:"this_month"`
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Tenant-ID")
		
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	tenantID := r.Header.Get("X-Tenant-ID")
	if tenantID == "" {
		http.Error(w, `{"error": "X-Tenant-ID header é obrigatório"}`, http.StatusBadRequest)
		return
	}

	var stats ProcessStats
	
	switch tenantID {
	case "11111111-1111-1111-1111-111111111111": // Silva & Associados
		stats = ProcessStats{Total: 45, Active: 38, Paused: 5, Archived: 2, ThisMonth: 12}
	case "22222222-2222-2222-2222-222222222222": // Costa Santos
		stats = ProcessStats{Total: 32, Active: 28, Paused: 3, Archived: 1, ThisMonth: 8}
	case "33333333-3333-3333-3333-333333333333": // Barros Empresa
		stats = ProcessStats{Total: 67, Active: 58, Paused: 7, Archived: 2, ThisMonth: 15}
	default:
		stats = ProcessStats{Total: 25, Active: 20, Paused: 3, Archived: 2, ThisMonth: 6}
	}

	log.Printf("📊 Stats para tenant %s: %+v", tenantID, stats)

	response := map[string]interface{}{
		"data":      stats,
		"timestamp": time.Now().UTC(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]interface{}{
		"status":    "healthy",
		"service":   "process-service-temp",
		"timestamp": time.Now().UTC(),
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/api/v1/processes/stats", statsHandler)

	handler := corsMiddleware(mux)

	log.Printf("🚀 Servidor iniciando na porta 8083")
	log.Printf("📊 Endpoint: GET /api/v1/processes/stats")
	log.Printf("🔍 Health: GET /health")

	if err := http.ListenAndServe(":8083", handler); err != nil {
		log.Fatal("Erro:", err)
	}
}