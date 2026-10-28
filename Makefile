.PHONY: all build test run lint clean docker-up docker-down

all: build test

build:
	mkdir -p bin
	go build -v -o bin/api ./cmd/api

test:
	go test -v -count=1 ./...

run:
	PORT=8080 go run ./cmd/api

lint:
	go vet ./...

clean:
	rm -rf bin
	go clean

docker-up:
	docker compose up -d

docker-down:
	docker compose down
