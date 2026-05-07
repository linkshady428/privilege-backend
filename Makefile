DATABASE_URL ?= postgres://postgres:postgres@localhost:5432/privilege?sslmode=disable

.PHONY: run up down logs psql seed migrate-up migrate-down migrate-new tidy

run:
	SKIP_AUTH=true PORT=8080 go run ./cmd/server

up:
	docker compose up --build

down:
	docker compose down -v

logs:
	docker compose logs -f server

psql:
	docker compose exec db psql -U postgres -d privilege

seed:
	docker compose exec -T db psql -U postgres -d privilege < scripts/seed.sql

migrate-up:
	docker compose run --rm migrate \
		-path=/migrations \
		-database postgres://postgres:postgres@db:5432/privilege?sslmode=disable \
		up

migrate-new:
	docker compose run --rm migrate \
		create -ext sql -dir /migrations -seq $(name)

migrate-down:
	docker compose run --rm migrate \
		-path=/migrations \
		-database postgres://postgres:postgres@db:5432/privilege?sslmode=disable \
		down 1

tidy:
	go mod tidy
