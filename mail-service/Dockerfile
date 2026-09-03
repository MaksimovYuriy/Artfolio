FROM golang:1.26-trixie AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /app/bin/mail-service ./cmd/mail

FROM debian:trixie-slim
RUN groupadd --system app \
    && useradd --system --gid app --no-create-home app
WORKDIR /app


COPY --from=build /app/bin/mail-service ./mail-service
USER app
ENTRYPOINT ["./mail-service"]