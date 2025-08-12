package portscan

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

var (
	currentScanner *Scanner
	isScanning     bool
)

// StartPortScanHandler обрабатывает запрос на начало сканирования портов
func StartPortScanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if isScanning {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Сканирование уже выполняется",
		})
		return
	}

	var config ScanConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Ошибка декодирования JSON: " + err.Error(),
		})
		return
	}

	// Валидация и нормализация конфигурации
	normalizedConfig, err := validateAndNormalizeScanConfig(config)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	// Создаем новый сканер
	currentScanner = NewScanner(normalizedConfig)
	isScanning = true

	// Запускаем сканирование в горутине
	go func() {
		defer func() {
			isScanning = false
		}()

		if err := currentScanner.Scan(); err != nil {
			fmt.Printf("Ошибка сканирования: %v\n", err)
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Сканирование портов запущено",
	})
}

// StopPortScanHandler обрабатывает запрос на остановку сканирования
func StopPortScanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	isScanning = false

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Сканирование остановлено",
	})
}

// PortScanStatsHandler возвращает статистику сканирования портов
func PortScanStatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]interface{}{
		"isScanning": isScanning,
		"results":    []PortResult{},
		"stats":      map[string]interface{}{"total": 0, "open": 0, "closed": 0},
	}

	if currentScanner != nil {
		response["results"] = currentScanner.GetResults()
		response["stats"] = currentScanner.GetStats()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// validateAndNormalizeScanConfig валидирует и нормализует конфигурацию сканирования
func validateAndNormalizeScanConfig(config ScanConfig) (ScanConfig, error) {
	// Проверяем target
	if config.Target == "" {
		return config, fmt.Errorf("не указан адрес для сканирования")
	}

	// Нормализуем target (убираем протокол если есть)
	target := config.Target
	if strings.HasPrefix(target, "http://") {
		target = strings.TrimPrefix(target, "http://")
	} else if strings.HasPrefix(target, "https://") {
		target = strings.TrimPrefix(target, "https://")
	}

	// Убираем путь если есть
	if strings.Contains(target, "/") {
		target = strings.Split(target, "/")[0]
	}

	// Убираем порт если есть
	if strings.Contains(target, ":") {
		target = strings.Split(target, ":")[0]
	}

	// Обновляем конфигурацию
	config.Target = target

	// Проверяем формат IP или домена
	if !isValidTarget(target) {
		return config, fmt.Errorf("неверный формат адреса: %s", target)
	}

	// Проверяем тип портов
	if config.PortType == "" {
		return config, fmt.Errorf("не указан тип портов для сканирования")
	}

	// Если custom, проверяем customPorts
	if config.PortType == "custom" {
		if config.CustomPorts == "" {
			return config, fmt.Errorf("не указаны пользовательские порты")
		}

		// Создаем временный сканер для валидации портов
		tempScanner := NewScanner(config)
		_, err := tempScanner.parseCustomPorts(config.CustomPorts)
		if err != nil {
			return config, fmt.Errorf("ошибка в пользовательских портах: %v", err)
		}
	} else {
		// Проверяем что тип портов существует
		if _, exists := portSets[config.PortType]; !exists {
			return config, fmt.Errorf("неизвестный тип портов: %s", config.PortType)
		}
	}

	// Проверяем таймаут
	if config.Timeout <= 0 {
		config.Timeout = 3 // значение по умолчанию
	}
	if config.Timeout > 30 {
		return config, fmt.Errorf("слишком большой таймаут (максимум 30 секунд): %d", config.Timeout)
	}

	return config, nil
}

// isValidTarget проверяет валидность цели (IP или домен)
func isValidTarget(target string) bool {
	// Проверяем IP адрес
	if isValidIP(target) {
		return true
	}

	// Проверяем доменное имя
	if isValidDomain(target) {
		return true
	}

	return false
}

// isValidIP проверяет валидность IP адреса
func isValidIP(ip string) bool {
	// Простая проверка IPv4
	ipRegex := regexp.MustCompile(`^(\d{1,3}\.){3}\d{1,3}$`)
	if !ipRegex.MatchString(ip) {
		return false
	}

	parts := strings.Split(ip, ".")
	for _, part := range parts {
		if len(part) > 3 {
			return false
		}

		var num int
		for _, char := range part {
			if char < '0' || char > '9' {
				return false
			}
			num = num*10 + int(char-'0')
		}

		if num > 255 {
			return false
		}
	}

	return true
}

// isValidDomain проверяет валидность доменного имени
func isValidDomain(domain string) bool {
	if len(domain) == 0 || len(domain) > 253 {
		return false
	}

	// Простая проверка доменного имени
	domainRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`)
	return domainRegex.MatchString(domain)
}

// GetPortTypesHandler возвращает доступные типы портов
func GetPortTypesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	portTypes := map[string]interface{}{
		"http": map[string]interface{}{
			"name":        "HTTP/HTTPS",
			"description": "Веб-серверы и HTTP сервисы",
			"ports":       portSets["http"],
		},
		"ssh": map[string]interface{}{
			"name":        "SSH",
			"description": "Secure Shell доступ",
			"ports":       portSets["ssh"],
		},
		"database": map[string]interface{}{
			"name":        "Базы данных",
			"description": "MySQL, PostgreSQL, MongoDB и другие",
			"ports":       portSets["database"],
		},
		"ftp": map[string]interface{}{
			"name":        "FTP",
			"description": "File Transfer Protocol",
			"ports":       portSets["ftp"],
		},
		"mail": map[string]interface{}{
			"name":        "Почтовые сервисы",
			"description": "SMTP, POP3, IMAP",
			"ports":       portSets["mail"],
		},
		"dns": map[string]interface{}{
			"name":        "DNS",
			"description": "Domain Name System",
			"ports":       portSets["dns"],
		},
		"gaming": map[string]interface{}{
			"name":        "Игровые сервисы",
			"description": "Minecraft, Steam и другие игры",
			"ports":       portSets["gaming"],
		},
		"common": map[string]interface{}{
			"name":        "Популярные порты",
			"description": "Наиболее часто используемые порты",
			"ports":       portSets["common"],
		},
		"custom": map[string]interface{}{
			"name":        "Пользовательские",
			"description": "Указать свои порты или диапазон",
			"ports":       []int{},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(portTypes)
}
