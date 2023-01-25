FROM golang:1.17-alpine AS build

WORKDIR /app

COPY . ./

RUN go mod download

RUN go build -o file-manager-app

FROM alpine:3.16

WORKDIR /app

COPY --from=build /app/file-manager-app .

EXPOSE 8080 8081

CMD ["./file-manager-app"]
