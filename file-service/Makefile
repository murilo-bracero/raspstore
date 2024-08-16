MAKEFLAGS += --silent

version := 0.0.1

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
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.26.0
	go install go.uber.org/mock/mockgen@v0.3.0

generate:
	sqlc generate -f config/sqlc.yaml
	mockgen -source=internal/application/repository/repository.go -destination=internal/application/repository/mocks/repository.go -package=mocks
	mockgen -source=internal/application/facade/file.go -destination=internal/application/facade/mocks/file.go -package=mocks
