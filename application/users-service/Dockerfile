FROM golang:1.17-alpine AS build

WORKDIR /app

COPY go.mod ./

COPY go.sum ./

RUN go mod download

COPY . ./

RUN go build -o users-service

FROM alpine:3.16

WORKDIR /app

COPY --from=build /app/users-service .

EXPOSE 9000 9090

CMD ["./users-service"]
