FROM golang:1.17-alpine as build

RUN apk add --no-cache protoc

WORKDIR /app

RUN GO111MODULE=on \
        go get google.golang.org/protobuf/cmd/protoc-gen-go \
        google.golang.org/grpc/cmd/protoc-gen-go-grpc

COPY authentication ./

COPY proto/authentication.proto ./

COPY proto/users-service.proto ./

RUN go mod download

RUN protoc --go_out=. --go-grpc_out=. *.proto

RUN go build -o authentication-app

FROM golang:1.17-alpine

WORKDIR /app

COPY --from=build /app/authentication-app authentication-app

EXPOSE 8080 8081

CMD ["/authentication-app"]
