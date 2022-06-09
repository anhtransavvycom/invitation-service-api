all: run

run:
	go run main.go

test:
	go test --cover ./...

start:
	docker compose up --build -d

stop:
	docker compose down

swagger_install:
	which swag || go install github.com/swaggo/swag/cmd/swag@latest

swagger: swagger_install
	swag init
