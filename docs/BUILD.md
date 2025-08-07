# 🚀 Load Tester - Инструкции по сборке

## 📋 Требования

- Go 1.22 или выше
- Git (для автоматического определения версии)

## 🛠️ Быстрая сборка

### Windows
```cmd
# Запустить скрипт сборки
build.bat

# Или использовать PowerShell
powershell -ExecutionPolicy Bypass -File build.ps1
```

### Linux/macOS
```bash
# Сделать скрипт исполняемым (только первый раз)
chmod +x build.sh

# Запустить сборку
./build.sh

# Или использовать Makefile
make all
```

## 🎯 Отдельные платформы

### Только Windows
```bash
make build-windows
# или
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o dist/windows/loadtester.exe ./cmd/loadtester
```

### Только Linux
```bash
make build-linux
# или
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o dist/linux/loadtester ./cmd/loadtester
```

### Только macOS
```bash
make build-macos
# или
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o dist/macos/loadtester ./cmd/loadtester
```

## 📦 Структура сборки

После выполнения сборки файлы будут размещены в следующей структуре:

```
dist/
├── windows/
│   └── loadtester-v1.0.0-windows-amd64.exe
├── linux/
│   ├── loadtester-v1.0.0-linux-amd64
│   └── loadtester-v1.0.0-linux-arm64
└── macos/
    ├── loadtester-v1.0.0-darwin-amd64    (Intel Mac)
    └── loadtester-v1.0.0-darwin-arm64    (Apple Silicon)
```

## 🏷️ Версионирование

Версия автоматически определяется из Git тегов:
- Если есть теги: используется последний тег
- Если тегов нет: используется хеш коммита
- Если Git недоступен: используется "v1.0.0"

### Создание релиза
```bash
# Создать тег версии
git tag v1.2.3

# Собрать с новой версией
make all

# Проверить версию
./dist/linux/loadtester-v1.2.3-linux-amd64 --version
```

## 🔧 Дополнительные команды

### Очистка
```bash
make clean
# или
rm -rf dist/
```

### Тестирование
```bash
make test
# или
go test -v ./...
```

### Локальный запуск
```bash
make run
# или
go run ./cmd/loadtester
```

### Информация о сборке
```bash
make info
```

### Создание архивов
```bash
make package
```

## 🎨 Флаги сборки

- `-s -w` - удаляет отладочную информацию (уменьшает размер)
- `-X main.version=...` - встраивает версию в бинарный файл
- `CGO_ENABLED=0` - отключает CGO для статической линковки

## 🚀 Оптимизация

### Минимальный размер
```bash
go build -ldflags "-s -w" -trimpath ./cmd/loadtester
```

### С компрессией UPX (если установлен)
```bash
go build -ldflags "-s -w" ./cmd/loadtester
upx --best loadtester
```

## 🐛 Устранение проблем

### Ошибка "command not found"
- Убедитесь, что Go установлен и добавлен в PATH
- Проверьте версию: `go version`

### Ошибка прав доступа (Linux/macOS)
```bash
chmod +x build.sh
chmod +x dist/linux/*
chmod +x dist/macos/*
```

### Проблемы с Git
```bash
# Если Git недоступен, установите версию вручную
VERSION=v1.0.0 make all
```

## 📋 Checklist релиза

- [ ] Обновить версию в git тегах
- [ ] Запустить полную сборку: `make all`
- [ ] Протестировать на каждой платформе
- [ ] Создать архивы: `make package`
- [ ] Проверить размеры файлов: `make info`
- [ ] Загрузить на GitHub Releases

## 🎯 Автоматизация CI/CD

Пример GitHub Actions workflow:

```yaml
name: Build and Release
on:
  push:
    tags: ['v*']
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-go@v4
      with:
        go-version: '1.22'
    - run: make all
    - run: make package
    - uses: softprops/action-gh-release@v1
      with:
        files: dist/*.zip,dist/*.tar.gz
```