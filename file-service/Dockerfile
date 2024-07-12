FROM golang:1.21-alpine AS build

WORKDIR /app

COPY --chown=185 . ./

RUN go mod download

RUN go build -o file-service ./cmd

FROM alpine:3.16

WORKDIR /app

COPY --chown=185 --from=build /app/file-service .

EXPOSE 8080 8080

CMD ["./file-service"]
