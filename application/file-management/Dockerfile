FROM golang:1.17-alpine

RUN apk add --no-cache protoc

WORKDIR /app

RUN GO111MODULE=on \
        go get google.golang.org/protobuf/cmd/protoc-gen-go \
        google.golang.org/grpc/cmd/protoc-gen-go-grpc

COPY go.mod ./

COPY go.sum ./

RUN go mod download

COPY . ./

RUN protoc --go_out=. --go-grpc_out=. ./messages/*.proto

RUN go build -o /file-manager-app

CMD ["/file-manager-app"]
