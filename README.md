# 🚀 Load Tester

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey)](https://github.com/715kg/loadtester/releases)

**Профессиональный инструмент для нагрузочного тестирования HTTP-серверов**

Load Tester - это мощный и простой в использовании инструмент для тестирования производительности веб-серверов и API. Программа создана для помощи разработчикам и системным администраторам в оценке отказоустойчивости их веб-приложений.

## ✨ Особенности

- 🌐 **Веб-интерфейс** - удобное управление через браузер
- ⚡ **Высокая производительность** - использует Go-рутины для максимальной эффективности  
- 📊 **Детальная статистика** - RPS, успешные запросы, ошибки, таймауты
- 🎯 **Кросс-платформенность** - Windows, Linux, macOS (Intel & Apple Silicon)
- 🛡️ **Безопасность** - встроенные ограничения и контроль ресурсов
- 📱 **Адаптивный дизайн** - работает на всех устройствах
- 🌙 **Темная/светлая тема** - переключение в реальном времени
- 📋 **Подробная документация** - инструкции и примеры использования

## 🚀 Быстрый старт

### Скачивание готовых сборок

1. Перейдите в [Releases](https://github.com/715kg/loadtester/releases)
2. Скачайте файл для вашей операционной системы:
   - **Windows:** `loadtester-v1.0.0-windows-amd64.exe`
   - **Linux AMD64:** `loadtester-v1.0.0-linux-amd64`
   - **Linux ARM64:** `loadtester-v1.0.0-linux-arm64`
   - **macOS Intel:** `loadtester-v1.0.0-darwin-amd64`
   - **macOS Apple Silicon:** `loadtester-v1.0.0-darwin-arm64`

### Запуск

#### Windows
```cmd
# Запустите скачанный файл
loadtester-v1.0.0-windows-amd64.exe

# Или проверьте версию
loadtester-v1.0.0-windows-amd64.exe --version
```

#### Linux/macOS
```bash
# Сделайте файл исполняемым
chmod +x loadtester-v1.0.0-linux-amd64

# Запустите программу
./loadtester-v1.0.0-linux-amd64

# Или проверьте версию
./loadtester-v1.0.0-linux-amd64 --version
```

### Использование

1. Откройте браузер и перейдите на http://localhost:8081
2. Настройте параметры тестирования
3. Примите пользовательское соглашение
4. Запустите тест и анализируйте результаты

## 🛠️ Сборка из исходного кода

### Требования

- Go 1.22 или выше
- Git (для автоматического определения версии)

### Клонирование репозитория

```bash
git clone https://github.com/715kg/loadtester.git
cd loadtester
```

### Быстрая сборка

#### Windows
```cmd
# CMD
scripts\build.bat

# PowerShell  
.\scripts\build.ps1
```

#### Linux/macOS
```bash
# Bash
chmod +x scripts/build.sh
./scripts/build.sh

# Make (если установлен)
make all
```

### Сборка для конкретной платформы

```bash
# Только для текущей платформы
go build -o loadtester ./cmd/loadtester

# Кросс-компиляция
make build-windows    # Windows
make build-linux      # Linux  
make build-macos      # macOS

# Или напрямую через Go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o loadtester.exe ./cmd/loadtester
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o loadtester ./cmd/loadtester
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o loadtester ./cmd/loadtester
```

### Запуск в режиме разработки

```bash
# Запуск
go run ./cmd/loadtester

# Проверка версии
go run ./cmd/loadtester --version

# Запуск тестов
go test ./...

# Запуск тестов с покрытием
go test -cover ./...

# Бенчмарки
go test -bench=. ./...
```

## 📁 Структура проекта

```
loadtester/
├── cmd/loadtester/          # Точка входа в программу
├── internal/                # Внутренние пакеты
│   ├── loadtest/           # Логика нагрузочного тестирования
│   └── templates/          # HTML шаблоны веб-интерфейса
├── pkg/info/               # Публичные пакеты (информация о программе)
├── docs/                   # Документация
├── scripts/                # Скрипты автоматизации сборки
├── assets/                 # Ресурсы и метаданные
└── dist/                   # Готовые бинарные файлы (создается при сборке)
```

## 🧪 Тестирование

```bash
# Запуск всех тестов
go test ./...

# Тесты с подробным выводом
go test -v ./...

# Тесты с покрытием кода
go test -cover ./...

# Только быстрые тесты (без интеграционных)
go test -short ./...

# Бенчмарки производительности
go test -bench=. ./...

# Тесты конкретного пакета
go test ./internal/loadtest
go test ./pkg/info
go test ./cmd/loadtester
```

## 📚 Документация

- **[📖 Инструкция по использованию](docs/BUILD.md)** - подробное руководство по сборке
- **[🛡️ Политика безопасности](docs/SECURITY.md)** - информация о безопасности
- **[🦠 Инструкции для антивирусов](docs/ANTIVIRUS.md)** - решение проблем с антивирусами
- **[📦 Инструкции для пользователей](docs/RELEASE.md)** - готовые сборки и установка
- **[📋 Информация о проекте](docs/PROJECT.md)** - детальная информация о проекте

## 🎯 Примеры использования

Подробные примеры конфигураций для разных типов тестирования см. в [инструкции по использованию](docs/RELEASE.md#примеры-использования).

## 📊 Системные требования

Подробные системные требования см. в [документации для пользователей](docs/RELEASE.md#системные-требования).

## 🤝 Вклад в проект

Мы приветствуем вклад в развитие проекта! Вы можете помочь:

- 🐛 Сообщать об ошибках через [Issues](https://github.com/715kg/loadtester/issues)
- 💡 Предлагать новые функции
- 📝 Улучшать документацию
- 🔧 Отправлять Pull Requests
- 🌍 Помогать с переводами

### Разработка

```bash
# Клонирование для разработки
git clone https://github.com/715kg/loadtester.git
cd loadtester

# Установка зависимостей
go mod download

# Запуск тестов
go test ./...

# Форматирование кода
go fmt ./...

# Проверка кода
go vet ./...
```

## ⚖️ Лицензия

Этот проект распространяется под лицензией MIT. Подробности в файле [LICENSE](LICENSE).

```
MIT License

Copyright (c) 2024 Алексей Ульянов

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.
```

## 👨‍💻 Автор

**Алексей Ульянов**
- GitHub: [@715kg](https://github.com/715kg)
- Проект: [loadtester](https://github.com/715kg/loadtester)

## 🛡️ Безопасность и ответственность

**⚠️ ВАЖНОЕ ПРЕДУПРЕЖДЕНИЕ:** Используйте программу только для тестирования собственных серверов или с явного разрешения владельца ресурса.

Подробная информация о безопасности и политике использования:
- **[🔒 Политика безопасности](docs/SECURITY.md)** - гарантии безопасности и процедуры
- **[🦠 Инструкции для антивирусов](docs/ANTIVIRUS.md)** - решение проблем с ложными срабатываниями

---

**⭐ Если проект оказался полезным, поставьте звездочку на GitHub!**