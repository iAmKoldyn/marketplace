.PHONY: build up down migrate

build:
    go build -o marketplace cmd/api/main.go

up: build
    docker-compose up --build

down:
    docker-compose down -v

build-worker:
    go build -o marketplace-worker cmd/worker/main.go
    
# run migrations locally
migrate:
    migrate -path migrations -database "$(POSTGRES_URL)" up
