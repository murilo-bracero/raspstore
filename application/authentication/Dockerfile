FROM golang:1.17-alpine as build

WORKDIR /app

COPY . ./

RUN go mod download

RUN go build -o authentication-app

FROM alpine:3.16

WORKDIR /app

COPY --from=build /app/authentication-app authentication-app

EXPOSE 8080 8081

CMD ["/authentication-app"]
