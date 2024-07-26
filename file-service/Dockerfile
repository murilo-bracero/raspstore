FROM golang:1.21-alpine AS build

RUN apk add --no-cache --update git build-base

WORKDIR /app

COPY --chown=185 . ./

RUN go mod download

RUN CGO_ENABLED=1 go build -o file-service ./cmd

FROM alpine:3.16

RUN addgroup -S rstoregroup && adduser -S rstoreuser -G rstoregroup

USER rstoreuser

WORKDIR /home/rstoreuser

COPY --chown=185 --from=build /app/file-service .
COPY --chown=185 --from=build /app/config ./config
COPY --chown=185 --from=build /app/db ./db

EXPOSE 8080 8080

CMD ["./file-service"]
