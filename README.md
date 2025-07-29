# Mobile App Config Service

## Описание
Сервис отдаёт конфигурацию мобильному приложению в зависимости от версии (`appVersion`) и платформы (`platform`).

Авторизация и проверка прав доступа не требуются.

### Структура пакетов:
- **`internal/app`** - основная структура приложения и конфигурация
- **`internal/storage`** - доступ к данным (база данных и репозитории)
- **`internal/service`** - бизнес-логика и резолвер
- **`internal/api`** - сгенерированные API хендлеры

## Запуск

### Запуск базы данных
```bash
cd deployments/db
docker-compose up -d
```

### Flow зависимостей:
```
Application -> Service (Resolver) -> Repositories -> Database
```
