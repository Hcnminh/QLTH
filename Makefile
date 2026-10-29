.PHONY: all deps build test test-browser qa run lint fmt-check docs clean docker-up docker-down db-init db-test

all: build test

deps:
	go mod download

build:
	mkdir -p bin
	go build -v -o bin/api ./cmd/api

test:
	go test -v -count=1 ./...

test-browser:
	NODE_PATH="$${NODE_PATH:-/usr/local/lib/node_modules}" node scripts/qa_browser.js

db-init:
	./scripts/init_db.sh

db-test:
	./scripts/test_db_schema.sh

qa:
	./scripts/qa_e2e.sh

run:
	PORT=$${PORT:-8098} go run ./cmd/api

lint:
	go vet ./...
	@out="$$(gofmt -l . 2>/dev/null)"; \
	if [ -n "$$out" ]; then echo "gofmt cần format lại:"; echo "$$out"; exit 1; fi
	@if command -v golangci-lint >/dev/null 2>&1; then golangci-lint run ./...; \
	elif [ -x "$$HOME/go/bin/golangci-lint" ]; then "$$HOME/go/bin/golangci-lint" run ./...; \
	else echo "golangci-lint chưa cài — bỏ qua (CI vẫn chạy đầy đủ)"; fi

docs:
	python3 scripts/validate_openapi.py
	python3 scripts/generate_api_docs.py

docs-check:
	python3 scripts/validate_openapi.py
	python3 scripts/generate_api_docs.py --check

clean:
	rm -rf bin
	go clean

docker-up:
	docker compose up -d

docker-down:
	docker compose down
