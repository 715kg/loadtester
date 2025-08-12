package portscan

import (
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

// PortResult представляет результат сканирования порта
type PortResult struct {
	Port        int    `json:"port"`
	IsOpen      bool   `json:"isOpen"`
	Service     string `json:"service"`
	Protocol    string `json:"protocol"`
	Description string `json:"description"`
}

// ScanConfig конфигурация для сканирования портов
type ScanConfig struct {
	Target      string `json:"target"`
	PortType    string `json:"portType"`    // "http", "ssh", "database", "custom", etc.
	CustomPorts string `json:"customPorts"` // "80,443" или "1-1000"
	Timeout     int    `json:"timeout"`     // в секундах
}

// Scanner структура для сканирования портов
type Scanner struct {
	config  ScanConfig
	results []PortResult
	mutex   sync.RWMutex
}

// Предопределенные наборы портов
var portSets = map[string][]int{
	"http":     {80, 443, 8080, 8443, 8000, 8888, 3000, 5000, 9000},
	"ssh":      {22, 2222},
	"database": {3306, 5432, 1433, 1521, 27017, 6379, 5984, 9042, 7000, 7001},
	"ftp":      {21, 22},
	"mail":     {25, 110, 143, 993, 995, 587, 465},
	"dns":      {53},
	"gaming":   {25565, 27015, 7777, 7778, 28960, 25575},
	"common":   {21, 22, 23, 25, 53, 80, 110, 143, 443, 993, 995, 3389, 5900},
}

// Описания сервисов для портов
var serviceDescriptions = map[int]string{
	21:    "FTP - File Transfer Protocol",
	22:    "SSH - Secure Shell",
	23:    "Telnet",
	25:    "SMTP - Simple Mail Transfer Protocol",
	53:    "DNS - Domain Name System",
	80:    "HTTP - HyperText Transfer Protocol",
	110:   "POP3 - Post Office Protocol v3",
	143:   "IMAP - Internet Message Access Protocol",
	443:   "HTTPS - HTTP Secure",
	993:   "IMAPS - IMAP over SSL",
	995:   "POP3S - POP3 over SSL",
	1433:  "MSSQL - Microsoft SQL Server",
	1521:  "Oracle Database",
	3306:  "MySQL Database",
	3389:  "RDP - Remote Desktop Protocol",
	5432:  "PostgreSQL Database",
	5900:  "VNC - Virtual Network Computing",
	6379:  "Redis Database",
	8080:  "HTTP Alternative",
	8443:  "HTTPS Alternative",
	25565: "Minecraft Server",
	27015: "Steam/Source Games",
	27017: "MongoDB Database",
}

// NewScanner создает новый сканер портов
func NewScanner(config ScanConfig) *Scanner {
	return &Scanner{
		config:  config,
		results: make([]PortResult, 0),
	}
}

// GetPortsToScan возвращает список портов для сканирования
func (s *Scanner) GetPortsToScan() ([]int, error) {
	if s.config.PortType == "custom" {
		return s.parseCustomPorts(s.config.CustomPorts)
	}

	if ports, exists := portSets[s.config.PortType]; exists {
		return ports, nil
	}

	return nil, fmt.Errorf("неизвестный тип портов: %s", s.config.PortType)
}

// parseCustomPorts парсит пользовательский ввод портов
func (s *Scanner) parseCustomPorts(input string) ([]int, error) {
	var ports []int

	parts := strings.Split(input, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)

		if strings.Contains(part, "-") {
			// Диапазон портов
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("неверный формат диапазона: %s", part)
			}

			start, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
			if err != nil {
				return nil, fmt.Errorf("неверный начальный порт: %s", rangeParts[0])
			}

			end, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
			if err != nil {
				return nil, fmt.Errorf("неверный конечный порт: %s", rangeParts[1])
			}

			if start > end {
				return nil, fmt.Errorf("начальный порт больше конечного: %d > %d", start, end)
			}

			if end-start > 1000 {
				return nil, fmt.Errorf("слишком большой диапазон портов (максимум 1000): %d-%d", start, end)
			}

			for i := start; i <= end; i++ {
				if i < 1 || i > 65535 {
					return nil, fmt.Errorf("порт вне допустимого диапазона (1-65535): %d", i)
				}
				ports = append(ports, i)
			}
		} else {
			// Отдельный порт
			port, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("неверный порт: %s", part)
			}

			if port < 1 || port > 65535 {
				return nil, fmt.Errorf("порт вне допустимого диапазона (1-65535): %d", port)
			}

			ports = append(ports, port)
		}
	}

	if len(ports) > 1000 {
		return nil, fmt.Errorf("слишком много портов для сканирования (максимум 1000): %d", len(ports))
	}

	// Удаляем дубликаты и сортируем
	portMap := make(map[int]bool)
	var uniquePorts []int
	for _, port := range ports {
		if !portMap[port] {
			portMap[port] = true
			uniquePorts = append(uniquePorts, port)
		}
	}

	sort.Ints(uniquePorts)
	return uniquePorts, nil
}

// ScanPort сканирует отдельный порт
func (s *Scanner) ScanPort(port int) PortResult {
	timeout := time.Duration(s.config.Timeout) * time.Second

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", s.config.Target, port), timeout)

	result := PortResult{
		Port:     port,
		IsOpen:   err == nil,
		Protocol: "TCP",
	}

	if err == nil {
		conn.Close()
		if desc, exists := serviceDescriptions[port]; exists {
			result.Description = desc
			// Извлекаем название сервиса из описания
			if idx := strings.Index(desc, " - "); idx != -1 {
				result.Service = desc[:idx]
			} else {
				result.Service = desc
			}
		} else {
			result.Service = "Unknown"
			result.Description = "Неизвестный сервис"
		}
	} else {
		result.Service = "Closed"
		result.Description = "Порт закрыт"
	}

	return result
}

// Scan выполняет сканирование портов
func (s *Scanner) Scan() error {
	ports, err := s.GetPortsToScan()
	if err != nil {
		return err
	}

	// Ограничиваем количество одновременных соединений
	maxWorkers := 50
	if len(ports) < maxWorkers {
		maxWorkers = len(ports)
	}

	portChan := make(chan int, len(ports))
	resultChan := make(chan PortResult, len(ports))

	// Запускаем воркеры
	var wg sync.WaitGroup
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range portChan {
				result := s.ScanPort(port)
				resultChan <- result
			}
		}()
	}

	// Отправляем порты в канал
	for _, port := range ports {
		portChan <- port
	}
	close(portChan)

	// Ждем завершения всех воркеров
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Собираем результаты
	s.mutex.Lock()
	s.results = make([]PortResult, 0, len(ports))
	for result := range resultChan {
		s.results = append(s.results, result)
	}
	s.mutex.Unlock()

	// Сортируем результаты по номеру порта
	sort.Slice(s.results, func(i, j int) bool {
		return s.results[i].Port < s.results[j].Port
	})

	return nil
}

// GetResults возвращает результаты сканирования
func (s *Scanner) GetResults() []PortResult {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	results := make([]PortResult, len(s.results))
	copy(results, s.results)
	return results
}

// GetOpenPorts возвращает только открытые порты
func (s *Scanner) GetOpenPorts() []PortResult {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var openPorts []PortResult
	for _, result := range s.results {
		if result.IsOpen {
			openPorts = append(openPorts, result)
		}
	}
	return openPorts
}

// GetStats возвращает статистику сканирования
func (s *Scanner) GetStats() map[string]interface{} {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	total := len(s.results)
	open := 0
	closed := 0

	for _, result := range s.results {
		if result.IsOpen {
			open++
		} else {
			closed++
		}
	}

	return map[string]interface{}{
		"total":  total,
		"open":   open,
		"closed": closed,
	}
}
