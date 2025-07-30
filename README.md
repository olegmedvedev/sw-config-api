# 🚀 Mobile App Config Service

## 📚 Документация

- **[📖 DEV_NOTES.md](docs/DEV_NOTES.md)** — детали реализации, архитектурные решения и технические особенности
- **[🧪 Tests README](tests/README.md)** — E2E тесты и сценарии тестирования API

---

## ⚡ Быстрый старт

### 🐙 Запуск всех сервисов (Docker Compose)

```bash
# Запуск приложения + база данных + Redis
docker-compose up -d
```

> **После запуска API будет доступно по адресу:** http://localhost:8080/config

### 🏗️ Запуск только инфраструктуры

```bash
# Только база данных и Redis
cd deployments/db
docker-compose up -d
```

### 📝 Переменные окружения

```bash
# Database configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=rootpassword
DB_NAME=sw_config

# Server configuration
SERVER_ADDR=:8080

# Redis configuration
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=
REDIS_DB=0
CACHE_TTL_SECONDS=300
```

### 🖥️ Локальная сборка

```bash
go build -o bin/sw-config-api ./cmd/sw-config-api
./bin/sw-config-api
```

### 🐳 Docker сборка

```bash
docker build -t sw-config-api .
docker run -p 8080:8080 sw-config-api
```

---

## 📝 Описание

Сервис отдаёт конфигурацию мобильному приложению в зависимости от версии (`appVersion`) и платформы (`platform`).

> **Авторизация и проверка прав доступа не требуются.**

### 🗂️ Структура пакетов

```text
internal/app      — основная структура приложения и конфигурация
internal/storage  — доступ к данным (база данных и репозитории)
internal/service  — бизнес-логика и резолвер
internal/api      — сгенерированные API хендлеры
internal/cache    — кэширование с Redis
```

---

## ✨ Возможности

- 🔄 Гибкая поддержка версий и платформ
- ⚡ Кэширование ответов через Redis
- 🧪 E2E и unit тесты
- 🐳 Простая интеграция с Docker и Docker Compose
- 🏗️ Graceful shutdown
- 📊 Структурированное логирование

---

## 🔗 Архитектура

```mermaid
flowchart TD
    A[Application] --> B[Service (Resolver)]
    B --> C[Repositories]
    C --> D[Database]
    B --> E[Cache (Redis)]
```

---

## 🧪 Тестирование

```bash
# Запуск E2E тестов
cd tests
go test -v

# С кастомным URL
E2E_BASE_URL=http://localhost:8080 go test -v
```
