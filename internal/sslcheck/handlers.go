package sslcheck

import (
	"encoding/json"
	"net/http"
)

var analyzer *SSLAnalyzer

func init() {
	analyzer = NewSSLAnalyzer()
}

type SSLCheckRequest struct {
	Target  string `json:"target"`
	Port    int    `json:"port"`
	Timeout int    `json:"timeout"`
}

type SSLCheckResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type SSLStatsResponse struct {
	IsAnalyzing bool        `json:"isAnalyzing"`
	Results     *SSLResults `json:"results"`
}

// StartSSLCheckHandler обрабатывает запрос на начало SSL анализа
func StartSSLCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req SSLCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := SSLCheckResponse{
			Success: false,
			Error:   "Ошибка декодирования JSON: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Валидация входных данных
	if req.Target == "" {
		response := SSLCheckResponse{
			Success: false,
			Error:   "Не указан целевой хост",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	if req.Port <= 0 || req.Port > 65535 {
		req.Port = 443 // Порт по умолчанию для HTTPS
	}

	if req.Timeout <= 0 {
		req.Timeout = 10 // Таймаут по умолчанию
	}

	// Запускаем анализ
	if err := analyzer.Start(req.Target, req.Port, req.Timeout); err != nil {
		response := SSLCheckResponse{
			Success: false,
			Error:   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := SSLCheckResponse{
		Success: true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// StopSSLCheckHandler обрабатывает запрос на остановку SSL анализа
func StopSSLCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	analyzer.Stop()

	response := SSLCheckResponse{
		Success: true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SSLStatsHandler возвращает текущие результаты SSL анализа
func SSLStatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := SSLStatsResponse{
		IsAnalyzing: analyzer.IsActive(),
		Results:     analyzer.GetResults(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RegisterSSLHandlers регистрирует все обработчики для SSL анализа
func RegisterSSLHandlers() {
	http.HandleFunc("/sslcheck/start", StartSSLCheckHandler)
	http.HandleFunc("/sslcheck/stop", StopSSLCheckHandler)
	http.HandleFunc("/sslcheck/stats", SSLStatsHandler)
}
