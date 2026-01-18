.PHONY: run

install:
	go mod tidy

build:
	go build -o build/server cmd/main.go

run:
	go run cmd/main.go

test:
	go test -v ./...

cover:
	go test -cover ./...

setup:
	go run .tools/rsa/keygen.go
	go run .tools/copy/envs.go 

# make apigen ARGS="sample"
apigen:
	go run .tools/apigen.go $(ARGS)

test-docker:
	@set -e; \
	docker compose --env-file .test.env -f docker-compose.test.yml build postgres-test redis-test; \
	docker compose --env-file .test.env -f docker-compose.test.yml up -d postgres-test redis-test; \
	docker compose --env-file .test.env -f docker-compose.test.yml run --rm migrate-test; \
	docker compose --env-file .test.env -f docker-compose.test.yml build goserver-test; \
	docker compose --env-file .test.env -f docker-compose.test.yml run --rm goserver-test; \
	docker compose --env-file .test.env -f docker-compose.test.yml down -v