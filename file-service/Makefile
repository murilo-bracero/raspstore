MAKEFLAGS += --silent

version := 0.0.1

golangci_version := v1.60.1
sqlc_version := v1.26.0
mockgen_version := v0.3.0

build: go-install go-build

go-build:
	go build -o file-service ./cmd

go-install:
	go mod download

docker-build:
	docker build \
		-t ghcr.io/murilo-bracero/spellscan-card-loader:$(version) \
		-t ghcr.io/murilo-bracero/spellscan-card-loader:latest \
		.

install-tools:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@$(sqlc_version)
	go install go.uber.org/mock/mockgen@$(mockgen_version)

generate:
	sqlc generate -f config/sqlc.yaml
	mockgen -source=internal/application/repository/repository.go -destination=internal/application/repository/mocks/repository.go -package=mocks
	mockgen -source=internal/application/facade/file.go -destination=internal/application/facade/mocks/file.go -package=mocks

go-lint:
	docker run -t --rm \
		-v $(shell pwd):/app \
		-w /app \
		golangci/golangci-lint:$(golangci_version) \
		golangci-lint run -v
