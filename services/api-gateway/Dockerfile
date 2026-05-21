# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS builder

RUN apk add --no-cache ca-certificates tzdata build-base git

WORKDIR /app

COPY . .

WORKDIR /app/services/api-gateway

RUN go mod tidy
RUN go build -o /out/api-gateway .

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata curl

WORKDIR /app

COPY --from=builder /out/api-gateway /app/api-gateway

ENV API_GATEWAY_PORT=8080

EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=5s --retries=5 CMD curl -f "http://localhost:8080/healthz" || exit 1

ENTRYPOINT ["/app/api-gateway"]
