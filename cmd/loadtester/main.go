package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/715kg/loadtester/internal/loadtest"
	"github.com/715kg/loadtester/internal/templates"
	"github.com/715kg/loadtester/pkg/info"
)

// Версия программы, устанавливается при сборке
var version = "dev"

// Дополнительная информация о программе
var (
	buildTime = "unknown"
	gitCommit = "unknown"
	buildUser = "unknown"
)

func main() {
	// Обрабатываем флаги командной строки
	showVersion := flag.Bool("version", false, "Показать версию программы")
	flag.Parse()

	if *showVersion {
		fmt.Printf("Load Tester %s\n", version)
		fmt.Printf("Автор: %s\n", info.Author)
		fmt.Printf("Лицензия: %s\n", info.License)
		fmt.Printf("Назначение: %s\n", info.Purpose)
		fmt.Printf("Категория: %s\n", info.Category)
		fmt.Printf("Сайт: %s\n", info.Website)
		if buildTime != "unknown" {
			fmt.Printf("Время сборки: %s\n", buildTime)
		}
		if gitCommit != "unknown" {
			fmt.Printf("Git коммит: %s\n", gitCommit)
		}
		fmt.Println("\nЭто легальный инструмент для тестирования производительности.")
		fmt.Println("Используйте только для собственных серверов или с разрешения владельца.")
		return
	}

	// Создаем HTTP сервер
	server := &http.Server{
		Addr:         ":8081",
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	http.HandleFunc("/", serveHTML)
	http.HandleFunc("/start", handleStart)
	http.HandleFunc("/stop", handleStop)
	http.HandleFunc("/stats", handleStats)

	// Канал для получения сигналов ОС
	sigChan := make(chan os.Signal, 1)

	// Регистрируем сигналы для graceful shutdown
	// SIGINT (Ctrl+C) - работает на всех платформах
	signal.Notify(sigChan, os.Interrupt)

	// Добавляем SIGTERM если доступен (Unix-like системы)
	signal.Notify(sigChan, syscall.SIGTERM)

	// Запускаем сервер в отдельной горутине
	go func() {
		fmt.Printf("🚀 %s %s запущен на http://localhost:8081\n", info.ProgramName, version)
		fmt.Printf("📝 Автор: %s | Лицензия: %s\n", info.Author, info.License)
		fmt.Printf("🎯 %s\n", info.ProgramDescription)
		fmt.Println("📖 Откройте браузер и перейдите по адресу выше")
		fmt.Println("⏹️  Для остановки нажмите Ctrl+C")
		fmt.Println("ℹ️  Это легальный инструмент для тестирования производительности")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Ошибка запуска сервера: %v\n", err)
			os.Exit(1)
		}
	}()

	// Ждем сигнал для завершения
	<-sigChan

	fmt.Println("\n🛑 Получен сигнал завершения, останавливаем сервер...")

	// Останавливаем текущий тест если он выполняется
	stopCurrentTest()

	// Создаем контекст с таймаутом для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Graceful shutdown сервера
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Ошибка при остановке сервера: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Сервер успешно остановлен")
}

func serveHTML(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("index").Parse(templates.HTMLTemplate))
	data := struct {
		Version string
	}{
		Version: version,
	}
	tmpl.Execute(w, data)
}

func handleStart(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var config loadtest.Config
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   "Ошибка парсинга конфигурации: " + err.Error(),
		})
		return
	}

	// Запускаем тест через функцию из loadtest пакета
	if err := loadtest.Start(&config); err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

func handleStop(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Останавливаем тест через функцию из loadtest пакета
	loadtest.Stop()

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
	})
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	// Получаем статистику через функцию из loadtest пакета
	stats := loadtest.GetStats()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// stopCurrentTest останавливает текущий выполняющийся тест
func stopCurrentTest() {
	if loadtest.IsRunning() {
		fmt.Println("⏹️  Останавливаем выполняющийся тест...")
		loadtest.Stop()

		// Ждем немного чтобы тест корректно завершился
		time.Sleep(2 * time.Second)

		fmt.Println("✅ Тест остановлен")
	}
}
