DATABASE_URL ?= postgres://postgres:postgres@localhost:5432/privilege?sslmode=disable

.PHONY: run up down migrate-up migrate-down migrate-new tidy

run:
	SKIP_AUTH=true PORT=8080 go run ./cmd/server

up:
	docker compose up --build

down:
	docker compose down -v

migrate-up:
	migrate -path migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path migrations -database "$(DATABASE_URL)" down 1

migrate-new:
	migrate create -ext sql -dir migrations -seq $(name)

tidy:
	go mod tidy
