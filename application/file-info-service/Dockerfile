FROM golang:1.17-alpine AS build

WORKDIR /app

COPY --chown=185 . ./

RUN go mod download

RUN go build -o file-info-service

FROM alpine:3.16

WORKDIR /app

COPY --chown=185 --from=build /app/file-info-service .

EXPOSE 8080 8081

CMD ["./file-info-service"]
