# Mobile App Config Service

## Описание
Сервис отдаёт конфигурацию мобильному приложению в зависимости от версии (`appVersion`) и платформы (`platform`).

Авторизация и проверка прав доступа не требуются.

### Структура пакетов:
- **`internal/app`** - основная структура приложения и конфигурация
- **`internal/storage`** - доступ к данным (база данных и репозитории)
- **`internal/service`** - бизнес-логика и резолвер
- **`internal/api`** - сгенерированные API хендлеры
- **`internal/cache`** - кэширование с Redis

## Запуск

### Запуск базы данных и Redis
```bash
cd deployments/db
docker-compose up -d
```

### Переменные окружения
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

### Сборка и запуск

#### Локальная сборка
```bash
# Сборка приложения
go build -o bin/sw-config-api ./cmd/sw-config-api

# Запуск
./bin/sw-config-api
```

#### Docker сборка
```bash
# Сборка образа
docker build -t sw-config-api .

# Запуск контейнера
docker run -p 8080:8080 sw-config-api
```

#### Запуск с Docker Compose
```bash
# Запуск всех сервисов (приложение + база + Redis)
docker-compose up -d

# Только база данных и Redis
cd deployments/db
docker-compose up -d
```

### Flow зависимостей:
```
Application -> Service (Resolver) -> Repositories -> Database
                    ↓
                Cache (Redis)
```
