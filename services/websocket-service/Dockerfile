# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS builder

RUN apk add --no-cache ca-certificates tzdata build-base git

WORKDIR /app

# Copy service and shared-go module files
COPY services/websocket-service/go.mod services/websocket-service/go.sum ./services/websocket-service/
COPY packages/shared-go/go.mod packages/shared-go/go.sum ./packages/shared-go/

WORKDIR /app/services/websocket-service

RUN go mod download

# Copy full repo so replace paths resolve
WORKDIR /app
COPY . .

WORKDIR /app/services/websocket-service
RUN go build -o /out/websocket-service ./...

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata curl

WORKDIR /app

COPY --from=builder /out/websocket-service /app/websocket-service

ENV WEBSOCKET_SERVICE_PORT=8084

EXPOSE 8084

HEALTHCHECK --interval=10s --timeout=5s --retries=5 CMD curl -f "http://localhost:8084/healthz" || exit 1

ENTRYPOINT ["/app/websocket-service"]
