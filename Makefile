include .env
export

export PROJECT_ROOT=${shell pwd}

env-up:
	@docker compose up -d sourcershub-postgres

env-down:
	@docker compose down sourcershub-postgres

env-cleanup:
	@read -p "Точно ли Вы хотите удалить все данные из БД? [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
	  docker compose down sourcershub-postgres port-forwarder && \
	  rm -rf out/pgdata && \
	  echo "БД очищена"; \
  	else \
  	  echo "Очистка БД отменена!"; \
  	fi

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Ошибка: укажите номер миграции."; \
		echo "Пример: make migrate-create seq=001"; \
		exit 1; \
	fi
	@docker compose run --rm sourcershub-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Ошибка: укажите действие."; \
		echo "Пример: make migrate-action action=up"; \
		exit 1; \
	fi
	@docker compose run --rm sourcershub-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@sourcershub-postgres:5432/${POSTGRES_DB}?sslmode=disable \
        ${action}

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

sourcershub-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run cmd/sourcers_hub/main.go

clear-logs:
	@test -n "$(PROJECT_ROOT)" || { echo "PROJECT_ROOT не задан!" >&2; exit 1; }
	@rm -rf ${PROJECT_ROOT}/out/logs && \
	echo "Логи успешно очищены!"