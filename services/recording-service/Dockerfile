FROM golang:1.26-alpine AS builder

RUN apk add --no-cache ca-certificates tzdata build-base git

WORKDIR /app
COPY . .

WORKDIR /app/services/recording-service
RUN go mod tidy
RUN go build -o /out/recording-service .

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata curl

WORKDIR /app
COPY --from=builder /out/recording-service /app/recording-service

ENV RECORDING_SERVICE_PORT=8083

EXPOSE 8083

HEALTHCHECK --interval=10s --timeout=5s --retries=5 CMD curl -f "http://localhost:8083/healthz" || exit 1

ENTRYPOINT ["/app/recording-service"]
