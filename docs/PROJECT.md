# 🚀 Load Tester - Информация о проекте

## 📋 Основная информация

- **Название:** Load Tester
- **Версия:** 1.0.0
- **Автор:** Алексей Ульянов
- **GitHub:** https://github.com/715kg/loadtester
- **Лицензия:** MIT License
- **Язык:** Go (Golang)

## 🎯 Описание

Load Tester - это профессиональный инструмент для нагрузочного тестирования веб-серверов и API. Программа создана для помощи разработчикам и системным администраторам в оценке производительности и отказоустойчивости их веб-приложений.

## ✨ Ключевые особенности

- 🌐 **Веб-интерфейс** - удобное управление через браузер
- ⚡ **Высокая производительность** - использует Go-рутины
- 📊 **Детальная статистика** - RPS, ошибки, таймауты
- 🛡️ **Безопасность** - встроенные ограничения
- 📱 **Адаптивный дизайн** - работает на всех устройствах
- 🌍 **Кросс-платформенность** - Windows, Linux, macOS

## 🏗️ Архитектура

```
Load Tester
├── cmd/loadtester/          # Точка входа в программу
├── internal/                # Внутренние пакеты
│   ├── loadtest/           # Логика нагрузочного тестирования
│   └── templates/          # HTML шаблоны веб-интерфейса
├── pkg/info/               # Публичные пакеты (информация о программе)
├── scripts/                # Скрипты автоматизации сборки
└── docs/                   # Документация
```

## 🛠️ Технологии

- **Backend:** Go 1.22+
- **Frontend:** HTML5, CSS3, JavaScript (Vanilla)
- **HTTP Client:** net/http с оптимизациями
- **Concurrency:** Go-рутины и каналы
- **Build:** Cross-compilation для всех платформ
- **Testing:** Unit, Integration и Benchmark тесты

## 📦 Поддерживаемые платформы

- **Windows:** x64 (.exe)
- **Linux:** x64, ARM64
- **macOS:** Intel, Apple Silicon

## 🔧 Сборка

```bash
# Быстрая сборка для всех платформ
make all

# Или используйте скрипты
./build.sh          # Linux/macOS
.\build.bat         # Windows CMD
.\build.ps1         # Windows PowerShell
```

## 📊 Статистика проекта

- **Размер исполняемого файла:** ~8-9 MB
- **Время сборки:** ~30 секунд
- **Поддерживаемые HTTP методы:** GET, POST, PUT, DELETE, HEAD, OPTIONS
- **Максимальная конкурентность:** 50,000+ горутин
- **Поддерживаемые протоколы:** HTTP/1.1, HTTP/2

## 🎨 Интерфейс

- **Главная страница:** Форма настройки тестирования
- **Инструкция:** Подробное руководство по использованию
- **Соглашение:** Правовая информация и ответственность
- **Темная/светлая тема:** Переключение в реальном времени
- **Прогресс-бар:** Визуализация выполнения теста
- **Статистика:** Детальные метрики в реальном времени

## 🛡️ Безопасность

- ✅ Открытый исходный код
- ✅ Не собирает личные данные
- ✅ Не модифицирует системные файлы
- ✅ Статическая компиляция без зависимостей
- ✅ Встроенная информация об авторе и назначении

## 📈 Производительность

### Рекомендуемые системные требования:
- **CPU:** 4+ ядра
- **RAM:** 8+ GB
- **Сеть:** 100+ Мбит/с
- **Диск:** 100 MB свободного места

### Максимальная производительность:
- **16+ ядер CPU:** До 50,000 одновременных соединений
- **32+ GB RAM:** До 500,000 запросов в тесте
- **1+ Гбит/с сеть:** До 10,000 RPS

## 📞 Поддержка

- **Issues:** https://github.com/715kg/loadtester/issues
- **Discussions:** https://github.com/715kg/loadtester/discussions

## 🤝 Вклад в проект

Приветствуются:
- Сообщения об ошибках
- Предложения улучшений
- Pull requests
- Документация
- Переводы

## 📄 Лицензия

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

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

## 🏆 Благодарности

Спасибо всем, кто использует и улучшает Load Tester!

---

**Контакты:**
- GitHub: https://github.com/715kg
- Проект: https://github.com/715kg/loadtester