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
		-t ghcr.io/murilo-bracero/raspstore-server:$(version) \
		-t ghcr.io/murilo-bracero/raspstore-server:latest \
		.

install-tools:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@$(sqlc_version)
	go install go.uber.org/mock/mockgen@$(mockgen_version)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin $(golangci_version)

generate:
	sqlc generate -f config/sqlc.yaml
	mockgen -source=internal/infra/repository/repository.go -destination=internal/infra/repository/mocks/repository.go -package=mocks
	mockgen -source=internal/application/facade/file.go -destination=internal/application/facade/mocks/file.go -package=mocks
	mockgen -source=internal/application/facade/fs.go -destination=internal/application/facade/mocks/fs.go -package=mocks

go-lint:
	golangci-lint run --enable gosec  -v
