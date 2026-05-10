# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o godis ./cmd/godis

# Run stage
FROM alpine:3.21

WORKDIR /app

COPY --from=builder /app/godis .

EXPOSE 6379

ENTRYPOINT ["./godis"]
