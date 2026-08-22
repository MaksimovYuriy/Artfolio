FROM golang:1.26-trixie AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /app/bin/artfolio-api ./cmd/api

FROM debian:trixie-slim
WORKDIR /app

COPY --from=build /app/bin/artfolio-api ./artfolio-api

EXPOSE 8081
CMD ["./artfolio-api"]