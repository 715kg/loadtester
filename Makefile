# Load Tester Makefile

# Получаем версию из git или используем дефолтную
VERSION ?= $(shell git describe --tags --always 2>/dev/null || echo "v1.0.0")

# Получаем дополнительную информацию для сборки
GIT_COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_USER ?= $(shell whoami 2>/dev/null || echo "unknown")
BUILD_TIME ?= $(shell date '+%Y-%m-%d %H:%M:%S' 2>/dev/null || echo "unknown")

# Расширенные флаги сборки
LDFLAGS = -ldflags "-s -w -X main.version=$(VERSION) -X main.gitCommit=$(GIT_COMMIT) -X 'main.buildTime=$(BUILD_TIME)' -X main.buildUser=$(BUILD_USER)"
BUILD_DIR = dist

# Цели по умолчанию
.PHONY: all clean build-windows build-linux build-macos test run help

all: clean build-windows build-linux build-macos

help:
	@echo "Load Tester Build System"
	@echo ""
	@echo "Доступные команды:"
	@echo "  all           - Собрать для всех платформ"
	@echo "  build-windows - Собрать для Windows"
	@echo "  build-linux   - Собрать для Linux"
	@echo "  build-macos   - Собрать для macOS"
	@echo "  test          - Запустить тесты"
	@echo "  run           - Запустить программу локально"
	@echo "  clean         - Очистить директорию сборки"
	@echo "  version       - Показать версию"

version:
	@echo "Версия: $(VERSION)"

clean:
	@echo "Очистка директории сборки..."
	@rm -rf $(BUILD_DIR)
	@echo "Готово!"

$(BUILD_DIR):
	@mkdir -p $(BUILD_DIR)/{windows,linux,macos}

build-windows: $(BUILD_DIR)
	@echo "Сборка для Windows (amd64)..."
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/windows/loadtester-$(VERSION)-windows-amd64.exe ./cmd/loadtester
	@echo "Windows сборка завершена!"

build-linux: $(BUILD_DIR)
	@echo "Сборка для Linux (amd64)..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/linux/loadtester-$(VERSION)-linux-amd64 ./cmd/loadtester
	@echo "Сборка для Linux (arm64)..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/linux/loadtester-$(VERSION)-linux-arm64 ./cmd/loadtester
	@chmod +x $(BUILD_DIR)/linux/*
	@echo "Linux сборка завершена!"

build-macos: $(BUILD_DIR)
	@echo "Сборка для macOS (amd64 - Intel)..."
	@CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/macos/loadtester-$(VERSION)-darwin-amd64 ./cmd/loadtester
	@echo "Сборка для macOS (arm64 - Apple Silicon)..."
	@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/macos/loadtester-$(VERSION)-darwin-arm64 ./cmd/loadtester
	@chmod +x $(BUILD_DIR)/macos/*
	@echo "macOS сборка завершена!"

test:
	@echo "Запуск тестов..."
	@go test -v ./...

run:
	@echo "Запуск Load Tester локально..."
	@go run ./cmd/loadtester

# Показать информацию о собранных файлах
info: $(BUILD_DIR)
	@echo "Информация о собранных файлах:"
	@echo ""
	@find $(BUILD_DIR) -type f -exec ls -lh {} \; | awk '{print $$9 " - " $$5}'
	@echo ""
	@echo "Общий размер:"
	@du -sh $(BUILD_DIR)

# Создать архивы для распространения
package: all
	@echo "Создание архивов для распространения..."
	@cd $(BUILD_DIR)/windows && zip -r ../loadtester-$(VERSION)-windows.zip *
	@cd $(BUILD_DIR)/linux && tar -czf ../loadtester-$(VERSION)-linux.tar.gz *
	@cd $(BUILD_DIR)/macos && tar -czf ../loadtester-$(VERSION)-macos.tar.gz *
	@echo "Архивы созданы в директории $(BUILD_DIR)/"