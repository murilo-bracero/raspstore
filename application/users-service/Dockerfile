FROM golang:1.17-alpine AS build

WORKDIR /app

COPY --chown=185 go.mod ./

COPY --chown=185 go.sum ./

RUN go mod download

COPY --chown=185 . ./

RUN go build -o users-service ./cmd

FROM alpine:3.16

WORKDIR /app

COPY --chown=185 --from=build /app/users-service .

EXPOSE 9000 9090

CMD ["./users-service"]
