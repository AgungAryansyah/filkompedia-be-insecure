# Build stage
FROM golang:1.23.6-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN apk add --no-cache make curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz
RUN mv migrate /usr/local/bin
RUN chmod +x /usr/local/bin/migrate

RUN go build -o main cmd/app/main.go

# Final stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate

CMD ["./main"]

