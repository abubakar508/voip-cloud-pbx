FROM golang:1.26-alpine AS builder

RUN apk add --no-cache ca-certificates tzdata build-base git

WORKDIR /app
COPY . .

WORKDIR /app/services/ai-service
RUN go mod tidy
RUN go build -o /out/ai-service .

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata curl

WORKDIR /app
COPY --from=builder /out/ai-service /app/ai-service

ENV AI_SERVICE_PORT=8085

EXPOSE 8085

HEALTHCHECK --interval=10s --timeout=5s --retries=5 CMD curl -f "http://localhost:8085/healthz" || exit 1

ENTRYPOINT ["/app/ai-service"]
