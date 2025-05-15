.PHONY: all build run stop clean proto test lint docker-up docker-down

# Переменные для сервисов
DOCKER_COMPOSE := docker compose
GO            := go
SERVICES      := films users sessions app
APP_PORT      := 8081
SESSIONS_PORT := 8010
FILMS_PORT    := 8020
USERS_PORT    := 8030

# Переменные для PostgreSQL
POSTGRES_USER     := postgres
POSTGRES_PASSWORD := postgres
POSTGRES_DB       := nimbus
POSTGRES_PORT     := 5432
POSTGRES_HOST     := localhost
MIGRATION_PATH    := /migrations

# Переменные для Redis
REDIS_PORT     := 6379
REDIS_HOST     := redis

# Порты для мониторинга
PROMETHEUS_PORT := 9090
GRAFANA_PORT   := 3000

# Цвета для вывода
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
RESET  := $(shell tput -Txterm sgr0)

# Основные команды
all: postgres-up redis-up migrate-up proto docker-build docker-up prometheus-up grafana-up

# Сборка всех сервисов локально
build:
	@echo "${GREEN}Сборка всех сервисов...${RESET}"
	$(GO) mod download
	@for service in $(SERVICES); do \
		echo "${YELLOW}Сборка $$service...${RESET}"; \
		$(GO) build -o ./bin/$$service ./cmd/$$service; \
	done

# Команды разработки без мониторинга
dev-up: postgres-up redis-up migrate-up docker-build docker-up
	@echo "${GREEN}Запуск проекта в режиме разработки (без мониторинга)...${RESET}"

dev-down: docker-down postgres-down redis-down
	@echo "${GREEN}Остановка сервисов разработки...${RESET}"

# Docker команды
docker-build:
	@echo "${GREEN}Сборка Docker образов...${RESET}"
	$(DOCKER_COMPOSE) build $(SERVICES)

docker-up:
	@echo "${GREEN}Запуск Docker контейнеров...${RESET}"
	$(DOCKER_COMPOSE) up -d $(SERVICES)

docker-down:
	@echo "${GREEN}Остановка Docker контейнеров...${RESET}"
	$(DOCKER_COMPOSE) down

# PostgreSQL команды
postgres-up:
	@echo "${GREEN}Запуск PostgreSQL...${RESET}"
	@docker-compose up -d postgres

postgres-down:
	@echo "${GREEN}Остановка PostgreSQL...${RESET}"
	@docker-compose stop postgres

postgres-create:
	@echo "${GREEN}Создание базы данных...${RESET}"
	@docker-compose exec -T postgres createdb --username=$(POSTGRES_USER) --owner=$(POSTGRES_USER) $(POSTGRES_DB)

postgres-drop:
	@echo "${GREEN}Удаление базы данных...${RESET}"
	@docker-compose exec -T postgres dropdb --username=$(POSTGRES_USER) $(POSTGRES_DB)

# Redis команды
redis-up:
	@echo "${GREEN}Запуск Redis...${RESET}"
	@docker-compose up -d redis

redis-down:
	@echo "${GREEN}Остановка Redis...${RESET}"
	@docker-compose stop redis

# Очистка и переинициализация базы
db-reinit: postgres-down
	@echo "${GREEN}Переинициализация базы данных...${RESET}"
	@docker-compose up -d postgres
	@sleep 3
	@docker-compose exec -T postgres dropdb --username=$(POSTGRES_USER) --if-exists $(POSTGRES_DB)
	@docker-compose exec -T postgres createdb --username=$(POSTGRES_USER) --owner=$(POSTGRES_USER) $(POSTGRES_DB)
	@docker-compose run --rm migrate \
		-path=$(MIGRATION_PATH) \
		-database="postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@postgres:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" \
		force $(shell docker-compose run --rm migrate \
			-path=$(MIGRATION_PATH) \
			-database="postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@postgres:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" \
			version)

# Очистка данных в PostgreSQL
clean-db-data:
	@echo "${GREEN}Очистка данных в базе данных...${RESET}"
	@docker-compose exec -T postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) -c "\
		DO \$$\$$ \
		DECLARE r RECORD; \
		BEGIN \
			FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') \
			LOOP \
				EXECUTE 'TRUNCATE TABLE ' || quote_ident(r.tablename) || ' CASCADE'; \
			END LOOP; \
		END \$$\$$;"

# Очистка данных в Redis
clean-redis-data:
	@echo "${GREEN}Очистка данных в Redis...${RESET}"
	@docker-compose exec -T redis redis-cli FLUSHALL

# Миграции
migrate-up:
	@echo "${GREEN}Применение миграций...${RESET}"
	@docker-compose run --rm migrate \
		-path=$(MIGRATION_PATH) \
		-database="postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@postgres:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" \
		up

migrate-down:
	@echo "${GREEN}Откат миграций...${RESET}"
	@docker-compose run --rm migrate \
		-path=$(MIGRATION_PATH) \
		-database="postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@postgres:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" \
		down

migrate-status:
	@echo "${GREEN}Статус миграций:${RESET}"
	@docker-compose run --rm migrate \
		-path=$(MIGRATION_PATH) \
		-database="postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@postgres:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" \
		version

# Мониторинг
prometheus-up:
	@echo "${GREEN}Запуск Prometheus...${RESET}"
	@docker-compose up -d prometheus
	@echo "${GREEN}Запуск cAdvisor...${RESET}"
	@docker-compose up -d cadvisor || echo "${YELLOW}Ошибка запуска cAdvisor, продолжаем...${RESET}"
	@echo "${GREEN}Запуск node_exporter...${RESET}"
	@docker-compose up -d node_exporter || echo "${YELLOW}Ошибка запуска node_exporter, продолжаем...${RESET}"

grafana-up:
	@echo "${GREEN}Запуск Grafana...${RESET}"
	@docker-compose up -d grafana

grafana-open:
	@echo "${GREEN}Открытие Grafana...${RESET}"
	@open http://localhost:$(GRAFANA_PORT)

prometheus-open:
	@echo "${GREEN}Открытие Prometheus...${RESET}"
	@open http://localhost:$(PROMETHEUS_PORT)

# Proto файлы
proto:
	@echo "${GREEN}Генерация proto файлов...${RESET}"
	cd proto && protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		./*.proto

# Swagger документация
swagger:
	@echo "${GREEN}Генерация документации Swagger для API Gateway...${RESET}"
	@mkdir -p docs/app
	@swag init -g cmd/app/main.go --parseDependency --output docs/app

# Разработка
lint:
	@echo "${GREEN}Запуск линтера...${RESET}"
	golangci-lint run ./...

test:
	@echo "${GREEN}Запуск тестов...${RESET}"
	$(GO) test -v -race -cover ./...

# Запуск локально
run: build
	@echo "${GREEN}Запуск всех сервисов...${RESET}"
	@./bin/sessions --b-port $(SESSIONS_PORT) & \
	./bin/films --b-port $(FILMS_PORT) & \
	./bin/users --b-port $(USERS_PORT) & \
	./bin/app --b-port $(APP_PORT)

# Остановка
stop:
	@echo "${GREEN}Остановка всех сервисов...${RESET}"
	@pkill -f "bin/(app|films|users|sessions)" || true

# Очистка
clean: postgres-down redis-down docker-down
	@echo "${GREEN}Очистка проекта...${RESET}"
	@rm -rf bin/
	docker system prune -f

# Статус системы
status:
	@echo "${GREEN}Статус Docker контейнеров:${RESET}"
	$(DOCKER_COMPOSE) ps
	@echo "\n${GREEN}Статус PostgreSQL:${RESET}"
	@docker-compose exec postgres pg_is_ready -U $(POSTGRES_USER) -d $(POSTGRES_DB) || echo "PostgreSQL не запущен"
	@echo "\n${GREEN}Статус Redis:${RESET}"
	@docker-compose exec redis redis-cli ping > /dev/null && echo "Redis ✅" || echo "Redis ❌"
	@echo "\n${GREEN}Статус метрик:${RESET}"
	@curl -s http://localhost:$(APP_PORT)/metrics > /dev/null && echo "API Gateway metrics ✅" || echo "API Gateway metrics ❌"
	@curl -s http://localhost:$(FILMS_PORT)/metrics > /dev/null && echo "Films service metrics ✅" || echo "Films service metrics ❌"
	@curl -s http://localhost:$(SESSIONS_PORT)/metrics > /dev/null && echo "Auth service metrics ✅" || echo "Auth service metrics ❌"
	@curl -s http://localhost:$(USERS_PORT)/metrics > /dev/null && echo "Users service metrics ✅" || echo "Users service metrics ❌"

# Помощь
help:
	@echo "${GREEN}Доступные команды:${RESET}"
	@echo "make all          - Запуск всего проекта"
	@echo "make build        - Локальная сборка всех сервисов"
	@echo "make docker-build - Сборка Docker образов"
	@echo "make docker-up    - Запуск в Docker"
	@echo "make docker-down  - Остановка Docker контейнеров"
	@echo "make postgres-up  - Запуск PostgreSQL"
	@echo "make redis-up     - Запуск Redis"
	@echo "make redis-down   - Остановка Redis"
	@echo "make migrate-up   - Применение миграций"
	@echo "make migrate-down - Откат миграций"
	@echo "make prometheus-up- Запуск Prometheus"
	@echo "make grafana-up   - Запуск Grafana"
	@echo "make grafana-open - Открыть Grafana в браузере"
	@echo "make run          - Локальный запуск всех сервисов"
	@echo "make stop         - Остановка локальных сервисов"
	@echo "make proto        - Генерация proto файлов"
	@echo "make swagger      - Генерация Swagger для API Gateway"
	@echo "make lint         - Запуск линтера"
	@echo "make test         - Запуск тестов"
	@echo "make status       - Проверка статуса всех компонентов"
	@echo "make clean        - Очистка проекта"
	@echo "make dev-up       - Запуск проекта без систем мониторинга"
	@echo "make dev-down     - Остановка сервисов разработки"
	@echo "make clean-redis-data - Очистка данных в Redis"