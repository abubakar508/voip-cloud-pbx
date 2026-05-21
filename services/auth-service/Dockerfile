# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS builder

RUN apk add --no-cache ca-certificates tzdata build-base git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /out/auth-service ./...

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata curl

WORKDIR /app

COPY --from=builder /out/auth-service /app/auth-service

ENV PORT=8081

EXPOSE 8081

HEALTHCHECK --interval=10s --timeout=5s --retries=5 CMD curl -f "http://localhost:8081/healthz" || exit 1

ENTRYPOINT ["/app/auth-service"]
