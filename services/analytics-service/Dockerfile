FROM golang:1.26-alpine AS builder

RUN apk add --no-cache ca-certificates tzdata build-base git

WORKDIR /app
COPY . .

WORKDIR /app/services/analytics-service
RUN go mod tidy
RUN go build -o /out/analytics-service .

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata curl

WORKDIR /app
COPY --from=builder /out/analytics-service /app/analytics-service

ENV ANALYTICS_SERVICE_PORT=8086

EXPOSE 8086

HEALTHCHECK --interval=10s --timeout=5s --retries=5 CMD curl -f "http://localhost:8086/healthz" || exit 1

ENTRYPOINT ["/app/analytics-service"]