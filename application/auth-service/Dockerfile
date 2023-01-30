FROM golang:1.17-alpine as build

WORKDIR /app

COPY --chown=185 . ./

RUN go mod download

RUN go build -o auth-service

FROM alpine:3.16

WORKDIR /app

COPY --chown=185 --from=build /app/auth-service auth-service

EXPOSE 9000 9090

CMD ["./auth-service"]
