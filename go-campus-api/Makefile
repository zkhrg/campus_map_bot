.PHONY: help run build test migration-up migration-down migration-create docker-up docker-down clean

help: ## Показать помощь
	@echo "Доступные команды:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

run: ## Запустить приложение локально
	go run cmd/app/main.go

build: ## Собрать приложение
	CGO_ENABLED=0 go build -o bin/app cmd/app/main.go

test: ## Запустить тесты
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

install-goose: ## Установить goose для миграций
	go install github.com/pressly/goose/v3/cmd/goose@latest

migration-up: ## Применить миграции
	goose -dir migrations postgres "$(DATABASE_URL)" up

migration-down: ## Откатить последнюю миграцию
	goose -dir migrations postgres "$(DATABASE_URL)" down

migration-status: ## Статус миграций
	goose -dir migrations postgres "$(DATABASE_URL)" status

migration-create: ## Создать новую миграцию (использование: make migration-create NAME=add_users_table)
	goose -dir migrations create $(NAME) sql

docker-up: ## Запустить Docker Compose
	docker-compose up --build -d

docker-down: ## Остановить Docker Compose
	docker-compose down

docker-logs: ## Показать логи Docker Compose
	docker-compose logs -f

docker-migrate: ## Применить миграции в Docker
	docker-compose exec app sh -c 'apk add --no-cache curl && curl -fsSL https://github.com/pressly/goose/releases/download/v3.15.0/goose_linux_x86_64 -o /usr/local/bin/goose && chmod +x /usr/local/bin/goose && goose -dir /migrations postgres "$(DATABASE_URL)" up'

clean: ## Очистить сборку
	rm -rf bin/
	rm -f coverage.out coverage.html

deps: ## Установить зависимости
	go mod download
	go mod tidy

lint: ## Запустить линтер
	golangci-lint run

fmt: ## Форматировать код
	go fmt ./...
	goimports -w .