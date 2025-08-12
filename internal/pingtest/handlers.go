package pingtest

import (
	"encoding/json"
	"net/http"
)

var pinger *PingTester

func init() {
	pinger = NewPingTester()
}

type PingTestRequest struct {
	Target   string `json:"target"`
	Count    int    `json:"count"`
	Interval int    `json:"interval"`
	Timeout  int    `json:"timeout"`
}

type PingTestResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type PingStatsResponse struct {
	IsRunning bool         `json:"isRunning"`
	Results   *PingResults `json:"results"`
}

// StartPingTestHandler обрабатывает запрос на начало ping теста
func StartPingTestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req PingTestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := PingTestResponse{
			Success: false,
			Error:   "Ошибка декодирования JSON: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Валидация входных данных
	if req.Target == "" {
		response := PingTestResponse{
			Success: false,
			Error:   "Не указан целевой хост",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	if req.Count <= 0 {
		req.Count = 4 // По умолчанию
	}
	if req.Count > 100 {
		req.Count = 100 // Максимум
	}

	if req.Interval <= 0 {
		req.Interval = 1 // По умолчанию 1 секунда
	}
	if req.Interval > 10 {
		req.Interval = 10 // Максимум 10 секунд
	}

	if req.Timeout <= 0 {
		req.Timeout = 3 // По умолчанию 3 секунды
	}
	if req.Timeout > 30 {
		req.Timeout = 30 // Максимум 30 секунд
	}

	// Запускаем ping тест
	if err := pinger.Start(req.Target, req.Count, req.Interval, req.Timeout); err != nil {
		response := PingTestResponse{
			Success: false,
			Error:   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := PingTestResponse{
		Success: true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// StopPingTestHandler обрабатывает запрос на остановку ping теста
func StopPingTestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pinger.Stop()

	response := PingTestResponse{
		Success: true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// PingStatsHandler возвращает текущие результаты ping теста
func PingStatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := PingStatsResponse{
		IsRunning: pinger.IsActive(),
		Results:   pinger.GetResults(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
