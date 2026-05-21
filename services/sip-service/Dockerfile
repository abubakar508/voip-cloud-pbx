FROM golang:1.26-alpine AS builder

RUN apk add --no-cache ca-certificates tzdata build-base git

WORKDIR /app

COPY . .

WORKDIR /app/services/sip-service

RUN go mod tidy
RUN go build -o /out/sip-service .

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata curl

WORKDIR /app

COPY --from=builder /out/sip-service /app/sip-service

ENV SIP_SERVICE_PORT=5070
ENV SIP_UDP_PORT=5060
ENV SIP_BIND_ADDR=0.0.0.0

EXPOSE 5070
EXPOSE 5060/udp

HEALTHCHECK --interval=10s --timeout=5s --retries=5 CMD curl -f "http://localhost:5070/healthz" || exit 1

ENTRYPOINT ["/app/sip-service"]
