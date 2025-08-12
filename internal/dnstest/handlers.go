package dnstest

import (
	"encoding/json"
	"net/http"
	"strings"
)

var dnsTester *DNSTester

func init() {
	dnsTester = NewDNSTester()
}

type DNSTestRequest struct {
	Domain      string   `json:"domain"`
	RecordTypes []string `json:"recordTypes"`
	DNSServers  []string `json:"dnsServers"`
	Timeout     int      `json:"timeout"`
}

type DNSTestResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}

type DNSStatsResponse struct {
	IsRunning bool        `json:"isRunning"`
	Results   *DNSResults `json:"results"`
}

type DNSTypesResponse struct {
	RecordTypes []RecordTypeInfo `json:"recordTypes"`
	DNSServers  []DNSServerInfo  `json:"dnsServers"`
}

type RecordTypeInfo struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type DNSServerInfo struct {
	IP          string `json:"ip"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// StartDNSTestHandler обрабатывает запрос на начало DNS теста
func StartDNSTestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req DNSTestRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := DNSTestResponse{
			Success: false,
			Error:   "Ошибка декодирования JSON: " + err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Валидация входных данных
	if req.Domain == "" {
		response := DNSTestResponse{
			Success: false,
			Error:   "Не указан домен для проверки",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Очищаем домен от протокола если есть
	domain := strings.TrimPrefix(req.Domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimPrefix(domain, "www.")
	req.Domain = domain

	// Валидация типов записей
	validTypes := map[string]bool{
		"A": true, "AAAA": true, "MX": true, "CNAME": true, "TXT": true, "NS": true,
	}

	if len(req.RecordTypes) == 0 {
		req.RecordTypes = []string{"A", "AAAA", "MX", "CNAME", "TXT", "NS"}
	} else {
		var validRecordTypes []string
		for _, recordType := range req.RecordTypes {
			if validTypes[recordType] {
				validRecordTypes = append(validRecordTypes, recordType)
			}
		}
		req.RecordTypes = validRecordTypes
	}

	if len(req.RecordTypes) == 0 {
		response := DNSTestResponse{
			Success: false,
			Error:   "Не указаны корректные типы DNS записей",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Валидация DNS серверов
	if len(req.DNSServers) == 0 {
		req.DNSServers = []string{"8.8.8.8", "1.1.1.1", "208.67.222.222"}
	}

	if req.Timeout <= 0 {
		req.Timeout = 5 // По умолчанию 5 секунд
	}
	if req.Timeout > 30 {
		req.Timeout = 30 // Максимум 30 секунд
	}

	// Запускаем DNS тест
	if err := dnsTester.Start(req.Domain, req.RecordTypes, req.DNSServers, req.Timeout); err != nil {
		response := DNSTestResponse{
			Success: false,
			Error:   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	response := DNSTestResponse{
		Success: true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// StopDNSTestHandler обрабатывает запрос на остановку DNS теста
func StopDNSTestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dnsTester.Stop()

	response := DNSTestResponse{
		Success: true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DNSStatsHandler возвращает текущие результаты DNS теста
func DNSStatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := DNSStatsResponse{
		IsRunning: dnsTester.IsActive(),
		Results:   dnsTester.GetResults(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetDNSTypesHandler возвращает доступные типы DNS записей и серверы
func GetDNSTypesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := DNSTypesResponse{
		RecordTypes: []RecordTypeInfo{
			{Type: "A", Name: "A записи", Description: "IPv4 адреса"},
			{Type: "AAAA", Name: "AAAA записи", Description: "IPv6 адреса"},
			{Type: "MX", Name: "MX записи", Description: "Почтовые серверы"},
			{Type: "CNAME", Name: "CNAME записи", Description: "Канонические имена"},
			{Type: "TXT", Name: "TXT записи", Description: "Текстовые записи"},
			{Type: "NS", Name: "NS записи", Description: "Серверы имен"},
		},
		DNSServers: []DNSServerInfo{
			{IP: "8.8.8.8", Name: "Google DNS", Description: "Публичный DNS Google"},
			{IP: "8.8.4.4", Name: "Google DNS", Description: "Альтернативный DNS Google"},
			{IP: "1.1.1.1", Name: "Cloudflare DNS", Description: "Публичный DNS Cloudflare"},
			{IP: "1.0.0.1", Name: "Cloudflare DNS", Description: "Альтернативный DNS Cloudflare"},
			{IP: "208.67.222.222", Name: "OpenDNS", Description: "Публичный DNS OpenDNS"},
			{IP: "208.67.220.220", Name: "OpenDNS", Description: "Альтернативный DNS OpenDNS"},
			{IP: "77.88.8.8", Name: "Yandex DNS", Description: "Публичный DNS Яндекс"},
			{IP: "77.88.8.1", Name: "Yandex DNS", Description: "Альтернативный DNS Яндекс"},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
