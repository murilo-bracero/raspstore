FROM golang:1.17-alpine AS build

WORKDIR /app

COPY --chown=185 . ./

RUN go mod download

RUN go build -o fs-service

FROM alpine:3.16

WORKDIR /app

COPY --chown=185 --from=build /app/fs-service .

EXPOSE 9000 9090

CMD ["./fs-service"]
