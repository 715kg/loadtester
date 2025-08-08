package loadtest

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"runtime"
	"runtime/debug"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Глобальные счетчики статистики
var (
	successCount uint64
	failCount    uint64
	errorCount   uint64

	// Детальная статистика по статус кодам
	status4xxCount  uint64 // 4xx ошибки клиента
	status5xxCount  uint64 // 5xx ошибки сервера
	timeoutCount    uint64 // Таймауты
	connectionCount uint64 // Ошибки соединения
	otherErrorCount uint64 // Прочие ошибки

	// Пул буферов что бы не создавать каждый раз объект
	bufferPool = sync.Pool{
		New: func() interface{} {
			buf := make([]byte, 4096)
			return &buf
		},
	}

	// Контроль выполнения теста
	testCancel        context.CancelFunc
	testRunning       bool
	testMutex         sync.RWMutex
	testStartTime     time.Time
	testTotalRequests int
)

// Config содержит конфигурацию для нагрузочного теста
type Config struct {
	TargetURL      string  `json:"targetURL"`
	Method         string  `json:"method"`
	TotalRequests  int     `json:"totalRequests"`
	MaxConcurrency int     `json:"maxConcurrency"`
	Timeout        int     `json:"timeout"` // в секундах
	MaxMemory      float64 `json:"maxMemory"`
	MaxCPUCores    int     `json:"maxCPUCores"` // количество ядер CPU, 0 = все доступные
}

// Stats содержит текущую статистику теста
type Stats struct {
	Total       uint64  `json:"total"`
	Success     uint64  `json:"success"`
	Fails       uint64  `json:"fails"`
	Errors      uint64  `json:"errors"`
	Status4xx   uint64  `json:"status4xx"`
	Status5xx   uint64  `json:"status5xx"`
	Timeouts    uint64  `json:"timeouts"`
	Connections uint64  `json:"connections"`
	RPS         float64 `json:"rps"`
	Elapsed     float64 `json:"elapsed"`
	IsRunning   bool    `json:"isRunning"`
	Progress    float64 `json:"progress"`
}

// Start запускает нагрузочный тест с заданной конфигурацией
func Start(config *Config) error {
	testMutex.Lock()
	defer testMutex.Unlock()

	if testRunning {
		return fmt.Errorf("тест уже выполняется")
	}

	// Валидация
	if config.TargetURL == "" {
		return fmt.Errorf("URL не может быть пустым")
	}

	// Сброс счетчиков
	resetCounters()

	testRunning = true
	testStartTime = time.Now()
	testTotalRequests = config.TotalRequests

	// Создаем контекст для отмены
	ctx, cancel := context.WithCancel(context.Background())
	testCancel = cancel

	// Запускаем тест в отдельной горутине
	go func() {
		defer func() {
			testMutex.Lock()
			testRunning = false
			testCancel = nil
			testMutex.Unlock()
		}()
		runLoadTestCore(ctx, config)
	}()

	return nil
}

// Stop останавливает выполняющийся тест
func Stop() {
	testMutex.Lock()
	defer testMutex.Unlock()

	if testCancel != nil {
		testCancel()
	}
}

// GetStats возвращает текущую статистику теста
func GetStats() Stats {
	testMutex.RLock()
	isRunning := testRunning
	startTime := testStartTime
	totalReq := testTotalRequests
	testMutex.RUnlock()

	success := atomic.LoadUint64(&successCount)
	fails := atomic.LoadUint64(&failCount)
	errors := atomic.LoadUint64(&errorCount)
	status4xx := atomic.LoadUint64(&status4xxCount)
	status5xx := atomic.LoadUint64(&status5xxCount)
	timeouts := atomic.LoadUint64(&timeoutCount)
	connections := atomic.LoadUint64(&connectionCount)
	total := success + fails + errors

	var elapsed float64
	var rps float64
	if !startTime.IsZero() {
		elapsed = time.Since(startTime).Seconds()
		if elapsed > 0 {
			rps = float64(total) / elapsed
		}
	}

	// Вычисляем прогресс
	var progress float64
	if totalReq > 0 {
		progress = float64(total) / float64(totalReq)
		if progress > 1 {
			progress = 1
		}
	}

	return Stats{
		Total:       total,
		Success:     success,
		Fails:       fails,
		Errors:      errors,
		Status4xx:   status4xx,
		Status5xx:   status5xx,
		Timeouts:    timeouts,
		Connections: connections,
		RPS:         rps,
		Elapsed:     elapsed,
		IsRunning:   isRunning,
		Progress:    progress,
	}
}

// IsRunning проверяет, выполняется ли тест в данный момент
func IsRunning() bool {
	testMutex.RLock()
	defer testMutex.RUnlock()
	return testRunning
}

// resetCounters сбрасывает все счетчики статистики
func resetCounters() {
	atomic.StoreUint64(&successCount, 0)
	atomic.StoreUint64(&failCount, 0)
	atomic.StoreUint64(&errorCount, 0)
	atomic.StoreUint64(&status4xxCount, 0)
	atomic.StoreUint64(&status5xxCount, 0)
	atomic.StoreUint64(&timeoutCount, 0)
	atomic.StoreUint64(&connectionCount, 0)
	atomic.StoreUint64(&otherErrorCount, 0)
}

// runLoadTestCore выполняет основную логику нагрузочного теста
func runLoadTestCore(ctx context.Context, config *Config) {
	// Настройка количества ядер CPU
	if config.MaxCPUCores > 0 {
		runtime.GOMAXPROCS(config.MaxCPUCores)
	}
	// Если MaxCPUCores = 0, используем все доступные ядра (поведение по умолчанию)

	// Настройка GC и памяти
	debug.SetGCPercent(50)
	debug.SetMemoryLimit(int64(GBToBytes(config.MaxMemory)))

	// Оптимизированный HTTP клиент
	client := &http.Client{
		Timeout: time.Duration(config.Timeout) * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:          config.MaxConcurrency * 2,
			MaxIdleConnsPerHost:   config.MaxConcurrency,
			MaxConnsPerHost:       config.MaxConcurrency,
			IdleConnTimeout:       30 * time.Second,
			TLSHandshakeTimeout:   5 * time.Second,
			ResponseHeaderTimeout: 5 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			DisableKeepAlives:     false,
			DisableCompression:    true,
			ForceAttemptHTTP2:     true,
		},
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, config.MaxConcurrency)

	// Запускаем нагрузку
requestLoop:
	for range config.TotalRequests {
		select {
		case <-ctx.Done():
			break requestLoop
		default:
			select {
			case semaphore <- struct{}{}:
				wg.Add(1)
				go func() {
					defer wg.Done()
					defer func() {
						<-semaphore
					}()
					makeRequest(ctx, client, config.TargetURL, config.Method)
				}()
			case <-ctx.Done():
				break requestLoop
			}
		}
	}

	// Ждём завершения всех горутин с таймаутом
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Все горутины завершены
	case <-time.After(10 * time.Second):
		// Таймаут ожидания
	}
}

// makeRequest выполняет HTTP запрос
func makeRequest(ctx context.Context, client *http.Client, targetURL, method string) {
	select {
	case <-ctx.Done():
		return
	default:
	}

	req, err := http.NewRequestWithContext(ctx, method, targetURL, nil)
	if err != nil {
		atomic.AddUint64(&errorCount, 1)
		return
	}

	// Добавляем заголовки
	req.Header.Set("User-Agent", "LoadTester/1.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Connection", "keep-alive")

	resp, err := client.Do(req)
	if err != nil {
		select {
		case <-ctx.Done():
			return
		default:
			classifyError(err)
		}
		return
	}
	defer resp.Body.Close()

	// Читаем и отбрасываем тело ответа
	if resp.Body != nil {
		bufPtr := bufferPool.Get().(*[]byte)
		defer bufferPool.Put(bufPtr)
		_, _ = io.CopyBuffer(io.Discard, io.LimitReader(resp.Body, 64*1024), *bufPtr)
	}

	// Классифицируем ответ по статус коду
	switch {
	case resp.StatusCode == http.StatusOK:
		atomic.AddUint64(&successCount, 1)
	case resp.StatusCode >= 400 && resp.StatusCode < 500:
		atomic.AddUint64(&status4xxCount, 1)
		atomic.AddUint64(&failCount, 1)
	case resp.StatusCode >= 500:
		atomic.AddUint64(&status5xxCount, 1)
		atomic.AddUint64(&failCount, 1)
	default:
		atomic.AddUint64(&failCount, 1)
	}
}

// classifyError классифицирует тип ошибки для детальной статистики
func classifyError(err error) {
	atomic.AddUint64(&errorCount, 1)

	errStr := err.Error()

	// Проверяем на таймауты
	if strings.Contains(errStr, "timeout") ||
		strings.Contains(errStr, "deadline exceeded") ||
		strings.Contains(errStr, "context deadline exceeded") {
		atomic.AddUint64(&timeoutCount, 1)
		return
	}

	// Проверяем на ошибки соединения
	if strings.Contains(errStr, "connection refused") ||
		strings.Contains(errStr, "connection reset") ||
		strings.Contains(errStr, "no such host") ||
		strings.Contains(errStr, "network is unreachable") {
		atomic.AddUint64(&connectionCount, 1)
		return
	}

	// Проверяем специфичные сетевые ошибки
	if netErr, ok := err.(net.Error); ok {
		if netErr.Timeout() {
			atomic.AddUint64(&timeoutCount, 1)
			return
		}
	}

	// Все остальные ошибки
	atomic.AddUint64(&otherErrorCount, 1)
}

// GBToBytes конвертирует гигабайты в байты
func GBToBytes(gb float64) uint64 {
	return uint64(gb * 1024 * 1024 * 1024)
}
