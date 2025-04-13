FROM golang:1.23.6-alpine

WORKDIR /app

COPY . .

RUN go mod tidy

RUN apk add --no-cache curl && \
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz && \
    mv migrate /usr/local/bin && \
    chmod +x /usr/local/bin/migrate

RUN apk add make

RUN go build -o main cmd/app/main.go

CMD ["./main"]
