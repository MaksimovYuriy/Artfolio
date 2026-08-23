FROM golang:1.26-trixie AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /app/bin/artfolio-api ./cmd/api
RUN go build -o /app/bin/artfolio-migrate ./cmd/migrate
RUN go build -o /app/bin/artfolio-key ./cmd/key

FROM debian:trixie-slim
WORKDIR /app

COPY --from=build /app/bin/artfolio-api ./artfolio-api
COPY --from=build /app/bin/artfolio-migrate ./artfolio-migrate
COPY --from=build /app/migrations ./migrations
COPY --from=build /app/bin/artfolio-key ./artfolio-key

EXPOSE 8081
CMD ["./artfolio-api"]