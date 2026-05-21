FROM golang:1.26-alpine AS builder

RUN apk add --no-cache ca-certificates tzdata build-base git

WORKDIR /app

# Copy entire repo so replace paths work
COPY . .

WORKDIR /app/services/media-service

RUN go mod tidy
RUN go build -o /out/media-service .

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata curl

WORKDIR /app

COPY --from=builder /out/media-service /app/media-service

ENV MEDIA_SERVICE_PORT=8082
ENV MEDIA_RTP_BASE_PORT=40000

EXPOSE 8082
EXPOSE 40000/udp

HEALTHCHECK --interval=10s --timeout=5s --retries=5 CMD curl -f "http://localhost:8082/healthz" || exit 1

ENTRYPOINT ["/app/media-service"]
