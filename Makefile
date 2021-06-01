.PHONY = migrate-up migrate-down dev dev-build

migrate = migrate -path db/migrations 
migrate += -database "postgres://$(user):$(password)@localhost:5432/$(db)?sslmode=disable"
compose = docker-compose

migrate-up:
	$(migrate) -verbose up

migrate-down:
	$(migrate) -verbose down

run:
	go run ./cmd/rodavis/main.go

dev:
	$(compose) -f docker-compose.dev.yml up

dev-build:
	$(compose) -f docker-compose.dev.yml up --build
  