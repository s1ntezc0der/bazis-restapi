# Task Management API

REST API сервис для управления задачами в командах с поддержкой ролевой модели, истории изменений и сложными SQL-запросами.

## Стек технологий

- **Go** — основной язык
- **MySQL** — база данных
- **Redis** — кеширование и rate limiting
- **Docker** — контейнеризация
- **Docker Compose** — оркестрация
- **Swagger** — документация API
- **Prometheus** — метрики
- **JWT** — аутентификация

## Быстрый старт

### 1. Клонировать репозиторий

```bash
	git clone https://github.com/s1ntezc0der/mkk_bazis.git
	cd mkk_bazis
```

### 2. Запустить через Docker Compose

```bash
	docker-compose up -d
```

### 3. Накатить миграции

```bash
	make migrate-up
```

### 4. Проверить, что сервис работает

```bash
	curl http://localhost:8080/api/v1/health
```

## API Эндпоинты

### Аутентификация

| Метод | Эндпоинт | Описание |
|-------|----------|----------|
| POST | `/api/v1/register` | Регистрация пользователя |
| POST | `/api/v1/login` | Вход (JWT) |

### Команды

| Метод | Эндпоинт | Описание |
|-------|----------|----------|
| POST | `/api/v1/teams` | Создать команду |
| GET | `/api/v1/teams` | Список команд пользователя |
| POST | `/api/v1/teams/{id}/invite` | Пригласить пользователя |

### Задачи

| Метод | Эндпоинт | Описание |
|-------|----------|----------|
| POST | `/api/v1/tasks` | Создать задачу |
| GET | `/api/v1/tasks` | Список задач (фильтрация + пагинация) |
| PUT | `/api/v1/tasks/{id}` | Обновить задачу |
| GET | `/api/v1/tasks/{id}/history` | История изменений |

### Комментарии

| Метод | Эндпоинт | Описание |
|-------|----------|----------|
| POST | `/api/v1/tasks/{id}/comments` | Добавить комментарий |
| GET | `/api/v1/tasks/{id}/comments` | Получить комментарии |

## Документация API (Swagger)

После запуска сервиса Swagger доступен по адресу:

```
http://localhost:8080/swagger/index.html
```

## Команды Make

```bash
	make migrate-up      # Применить миграции
	make migrate-down    # Откатить миграции
	make migrate-create  # Создать новую миграцию (make migrate-create name=имя)
	make test            # Запустить тесты
	make run             # Запустить сервер локально
	make docker-up       # Поднять контейнеры
	make docker-down     # Остановить контейнеры
```

## Переменные окружения (`.env`)

```env
	# База данных
	DB_USER=bazis
	DB_PASSWORD=bazis123
	DB_HOST=mysql
	DB_PORT=3306
	DB_NAME=bazis

	# Сервер
	PORT=8080

	# JWT
	JWT_SECRET=supersecret
	JWT_EXPIRE=24

	# Redis
	REDIS_ADDR=redis:6379

	# Prometheus
	METRICS_ENABLED=true
```

## Структура проекта

```
.
├── cmd/
│   ├── main.go          # Точка входа
│   └── migrate/
│       └── main.go      # Утилита миграций
├── internal/
│   ├── config/          # Конфигурация
│   ├── run/             # Запуск сервера
│   └── services/        # Бизнес-логика
│       ├── auth/        # Авторизация
│       ├── teams/       # Команды
│       └── tasks/       # Задачи + комментарии
├── pkg/
│   ├── cache/           # Redis-кеш
│   ├── db/              # MySQL
│   ├── errors/          # Кастомные ошибки
│   ├── jwt/             # JWT
│   ├── logger/          # Логи
│   ├── metrics/         # Prometheus
│   └── middleware/      # Middleware (auth, rate limit, circuit breaker)
├── migrations/          # SQL-миграции
├── tests/               # Unit и интеграционные тесты
├── docs/                # Swagger документация
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── .env
└── README.md
```

## Метрики (Prometheus)

Метрики доступны по адресу:

```
http://localhost:8080/metrics
```

## Тестирование

```bash
# Unit-тесты
go test ./tests/...

# С покрытием
go test ./tests/... -cover

# Интеграционные тесты
go test ./tests/integration/... -tags=integration
```
