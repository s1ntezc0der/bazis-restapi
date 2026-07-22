.PHONY: help migrate-up migrate-down migrate-create run test docker-up docker-down

help:
	@echo "Доступные команды:"
	@echo "  make migrate-up      - Применить миграции"
	@echo "  make migrate-down    - Откатить миграции"
	@echo "  make migrate-create  - Создать миграцию (make migrate-create name=имя)"
	@echo "  make run             - Запустить сервер локально"
	@echo "  make test            - Запустить тесты"
	@echo "  make docker-up       - Поднять контейнеры"
	@echo "  make docker-down     - Остановить контейнеры"

migrate-up:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down

migrate-create:
	go run cmd/migrate/main.go create $(name)

run:
	go run cmd/main.go

test:
	go test ./tests/... -v -cover

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

swagger:
	swag init -g cmd/main.go -o docs/

.PHONY: all
all: docker-up migrate-up
	@echo "✅ Сервис запущен и миграции применены"