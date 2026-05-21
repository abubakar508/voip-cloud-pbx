# VoIP Cloud PBX – Monorepo

A full VoIP Cloud PBX platform built with Go microservices and Next.js frontends.


 multi-service VoIP Cloud PBX platform built with Go (backend microservices) and Next.js (frontends). It includes:

- Auth service with JWT and multi-tenant users
- API gateway with JWT validation and reverse proxying
- SIP service for SIP signaling (REGISTER/INVITE/BYE/CANCEL)
- Media service for RTP handling, QoS tracking, and basic RTP forwarding
- WebSocket service for real-time events
- Dashboard web app for operators
- Shared Go and TypeScript packages

---

## High-Level Architecture

**Core components:**

- **Auth Service (`services/auth-service`)**  
  Provides tenant and user authentication:
  - `POST /auth/register` – create tenant + tenant admin user  
  - `POST /auth/login` – issue access + refresh JWTs  
  - `POST /auth/refresh` – refresh access token  
  - `GET /auth/me` – current user info (Bearer token)  

- **API Gateway (`services/api-gateway`)**  
  Single entry point for HTTP clients:
  - Validates JWT access tokens
  - Proxies `/auth/*` to auth-service
  - Exposes protected `/api/*` routes for future services

- **WebSocket Service (`services/websocket-service`)**  
  Real-time channel for frontend:
  - `/ws` endpoint (JWT-authenticated)
  - Tracks connected clients and supports broadcasting

- **SIP Service (`services/sip-service`)**  
  SIP signaling microservice:
  - Listens on UDP 5060 for SIP messages
  - Handles REGISTER:
    - Parses `To` and `Contact`
    - Looks up SIP account in Postgres
    - Stores registration binding in Redis (with TTL)
  - Handles INVITE/BYE/CANCEL:
    - Maintains in-memory call sessions
    - Publishes `calls.started` and `calls.ended` events to NATS

- **Media Service (`services/media-service`)**  
  RTP and media pipeline:
  - Listens for RTP on UDP port (default: 40000)
  - Parses RTP packets with Pion RTP
  - Tracks QoS per SSRC/address (packet count, loss, last seq)
  - Forwards RTP between two configured endpoints (legs A and B) per call
  - Subscribes to `calls.started` / `calls.ended` from NATS
  - HTTP endpoints:
    - `GET /healthz` – health check
    - `GET /calls` – current/ended call sessions
    - `GET /qos` – QoS snapshot
    - `POST /calls/:callId/endpoints` – configure RTP endpoints for a call

- **Dashboard Web App (`apps/dashboard-web`)**  
  Operator dashboard (Next.js):
  - Login via API gateway `/auth/login`
  - Shows current user and tenant
  - Displays:
    - Active/ended calls (from media-service `/calls`)
    - QoS stats (from media-service `/qos`)

- **Shared Libraries**
  - `packages/shared-go` – common Go utilities:
    - Config loading (Postgres, Redis, NATS, JWT, etc.)
    - Zap logger
    - HTTP server wrapper (Gin + `/healthz`)
  - `packages/shared-ts` – shared TypeScript types for frontend apps


---


Most configuration is via environment variables. Common ones:

- **Postgres**
  - `POSTGRES_HOST` (default: `postgres`)
  - `POSTGRES_PORT` (default: `5432`)
  - `POSTGRES_USER` (default: `voip`)
  - `POSTGRES_PASSWORD` (default: `voip_password`)
  - `POSTGRES_DB` (default: `voip_cloud_pbx`)

- **Redis**
  - `REDIS_HOST` (default: `redis`)
  - `REDIS_PORT` (default: `6379`)
  - `REDIS_PASSWORD` (optional)

- **NATS**
  - `NATS_URL` (default: `nats://nats:4222`)

- **JWT**
  - `JWT_ACCESS_SECRET`
  - `JWT_REFRESH_SECRET`
  - `JWT_ACCESS_TTL` (e.g., `15m`)
  - `JWT_REFRESH_TTL` (e.g., `720h`)

- **Ports**
  - Auth service: `8081`
  - API gateway: `8080`
  - WebSocket service: `8084`
  - SIP HTTP: `5070`, SIP UDP: `5060`
  - Media HTTP: `8082`, RTP UDP: `40000`
  - Dashboard: e.g., `3000` (depending on your Next.js dev/prod config)
  - Traefik: as configured in `docker-compose.yml`

---

This README walks a new user from:

1. Cloning the repo  
2. Starting all infrastructure and services  
3. Logging into the dashboard  
4. Hitting core endpoints and customizing configuration

---

## 1. Clone the repository

```bash
git clone https://github.com/abubakar508/voip-cloud-pbx.git
cd voip-cloud-pbx
```

The important top-level folders are:

- `infrastructure/` – Docker Compose, infra helpers, migrations
- `services/` – Go microservices
- `apps/` – Frontend apps (Next.js)
- `packages/` – Shared Go and TypeScript libraries

---

## 2. Prerequisites

Install:

- Docker and Docker Compose
- Git
- (Optional, for local development) Go 1.21+ and Node.js 18+

You do **not** need Go/Node to run the full stack via Docker; they’re only required if you want to run services or frontends directly on your machine.

---

## 3. Configure infrastructure (local)

All Docker orchestration for local runs is inside `infrastructure/`.

### 3.1 – Copy and adjust env file

From the `infrastructure/` folder:

```bash
cd infrastructure
cp local.env.example .env
```

Open `.env` and adjust if needed:

```env
APP_ENV=development
LOG_LEVEL=debug

POSTGRES_USER=voip
POSTGRES_PASSWORD=voip_password
POSTGRES_DB=voip_cloud_pbx

JWT_ACCESS_SECRET=dev_access_secret_change_me
JWT_REFRESH_SECRET=dev_refresh_secret_change_me
JWT_ACCESS_TTL=15m
JWT_REFRESH_TTL=720h
```

For a first run, you can keep defaults and just ensure secrets are non-empty strings.

---

## 4. Start the full stack with Docker

From `infrastructure/`:

```bash
# Build all images (first time or when code changes)
docker compose --env-file .env build

# Start everything in the background
docker compose --env-file .env up -d

# See status
docker compose --env-file .env ps
```

Key containers that should be running:

- `voip-postgres` – Postgres
- `voip-redis` – Redis
- `voip-nats` – NATS
- `voip-traefik` – Traefik reverse proxy
- `voip-auth-service` – Auth microservice
- `voip-api-gateway` – API gateway
- `voip-websocket-service` – WebSocket service
- `voip-sip-service` – SIP service (UDP 5060)
- `voip-media-service` – Media/RTP service (UDP 40000)
- `voip-recording-service` – Recording service skeleton
- `voip-ai-service` – AI service skeleton
- `voip-analytics-service` – Analytics service
- `voip-dashboard-web` – Dashboard frontend

To tail logs:

```bash
docker compose --env-file .env logs -f
```

---

## 5. Services and URLs

Traefik exposes HTTP services by hostname on port 80.

Add these hostnames to your local `/etc/hosts` (or equivalent) if you want to use them:

```text
127.0.0.1  api.localhost
127.0.0.1  dashboard.localhost
127.0.0.1  media.localhost
127.0.0.1  sip.localhost
127.0.0.1  ws.localhost
127.0.0.1  recording.localhost
127.0.0.1  ai.localhost
127.0.0.1  analytics.localhost
```

Then you can access:

- API gateway: `http://api.localhost`
- Dashboard: `http://dashboard.localhost`
- Media service: `http://media.localhost`
- SIP HTTP health: `http://sip.localhost/healthz`
- Traefik dashboard: `http://localhost:8080` (optional; can be disabled later)

If you don’t want to edit hosts, you can still use `localhost` with `Host` header overrides in curl when testing.

---

## 6. First login flow (Auth + Dashboard)

### 6.1 – Register a tenant admin

Call the API gateway (which proxies to auth-service):

```bash
curl -X POST http://api.localhost/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "tenantName": "Demo Tenant",
    "email": "admin@demo.local",
    "password": "Passw0rd!",
    "displayName": "Demo Admin"
  }'
```

Expected response: JSON containing `userId` and `tenantId`.

### 6.2 – Login

```bash
curl -X POST http://api.localhost/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@demo.local",
    "password": "Passw0rd!"
  }'
```

You should get:

- `accessToken`
- `refreshToken`
- `userId`
- `tenantId`

### 6.3 – Open the dashboard

In your browser:

- Go to `http://dashboard.localhost`
- Enter the same email and password (`admin@demo.local` / `Passw0rd!`)

After login, the dashboard should show:

- Logged-in user and tenant
- Tables for:
  - Active/ended calls (from media-service `/calls`)
  - QoS streams (from media-service `/qos`)

(Initially these may be empty until you have SIP/RTP traffic.)

---

## 7. SIP and media: basic usage

### 7.1 – SIP registration

1. Ensure `voip-sip-service` is running and port `5060/udp` is mapped (compose does this).
2. Create SIP accounts in Postgres (in `sip_accounts` table) with bcrypt-hashed passwords.
3. Configure a SIP client (Linphone, Zoiper, etc.):

   - SIP server: your machine IP (or `localhost` if calling from host)
   - Port: 5060
   - Transport: UDP
   - Username: SIP account username
   - Password: SIP account password

When the client sends REGISTER:

- `sip-service` parses the SIP REGISTER.
- Validates that the `Authorization` header’s username matches a known `sip_accounts` entry.
- Stores registration binding in Redis with TTL.
- Responds `200 OK` on success.

### 7.2 – SIP calls and call events

When placing a call (INVITE):

- `sip-service` creates an in-memory call session.
- Publishes:
  - `calls.started` when the call begins.
  - `calls.ended` on BYE/CANCEL.
- `media-service`, `recording-service`, `ai-service`, and `analytics-service` all consume these events via NATS.

You can observe:

- `media-service` logs `registered call session from NATS (started)` and `updated call session (ended)`.
- `analytics-service` writes basic `CallAnalytics` records to Postgres.

---

## 8. Media / RTP / QoS / Forwarding

### 8.1 – RTP listener

`media-service` listens on UDP port `40000` and parses RTP packets:

- Logs SSRC, sequence number, timestamp, payload type, payload length.
- Updates QoS stats per `(SSRC, remote address)`.

### 8.2 – QoS endpoint

Visit:

```bash
curl http://media.localhost/qos
```

You’ll get a JSON array like:

```json
[
  {
    "key": { "ssrc": 12345678, "addr": "127.0.0.1:54321" },
    "stats": { "packets": 100, "lost": 2, "lastSeq": 3456 }
  }
]
```

This is what the dashboard uses for the QoS table.

### 8.3 – Calls endpoint

```bash
curl http://media.localhost/calls
```

Returns current call sessions known by media-service, populated from NATS `calls.started` / `calls.ended`.

### 8.4 – RTP forwarding test

`media-service` can forward RTP between two endpoints (legs A and B) for a call.

1. Assume you have a `callId` (from SIP logs or just a test ID).
2. Set endpoints:

```bash
curl -X POST http://media.localhost/calls/<callId>/endpoints \
  -H "Content-Type: application/json" \
  -d '{"aAddr":"127.0.0.1:5000","bAddr":"127.0.0.1:5002"}'
```

3. Send RTP to `media-service` (UDP 40000) **from** port 5000 and listen on port 5002.  
   RTP arriving from `127.0.0.1:5000` will be forwarded to `127.0.0.1:5002`.

This is a simple test harness for your media pipeline.

---

## 9. Recording, AI, Analytics

These services are wired and ready but intentionally minimal.

- **recording-service**
  - Listens to `calls.started` and `calls.ended` via NATS.
  - Has `Recording` and `CallRecord` models in Postgres.
  - Intended to handle recording files from media-service in future.

- **ai-service**
  - Listens to `calls.ended`.
  - Has an `AISummary` model in Postgres.
  - Intended to generate summaries, sentiment, and keywords later.

- **analytics-service**
  - Listens to `calls.started` and `calls.ended`.
  - Writes simple `CallAnalytics` rows (start time, end time, duration) to Postgres.
  - Can be extended with aggregated reports.

You can open connections to Postgres and inspect these tables as calls are made.

---

## 10. Customizing for your own environment

### 10.1 – Change ports, hosts, secrets

- Edit `infrastructure/.env`:
  - Change `JWT_*` secrets for production-like setups.
- Edit `infrastructure/docker-compose.yml`:
  - Change published ports.
  - Change Traefik host rules:
    - `api.localhost`, `dashboard.localhost`, etc.
  - Add TLS or custom domains as needed.

### 10.2 – Running a single service locally

You can run a service outside Docker for development while still using Docker for infra (Postgres/Redis/NATS):

Example for auth-service:

```bash
# Ensure infra is up
cd infrastructure
docker compose --env-file .env up -d postgres redis nats

# In another terminal
cd ../services/auth-service

APP_ENV=development \
LOG_LEVEL=debug \
POSTGRES_HOST=localhost \
POSTGRES_PORT=5432 \
POSTGRES_USER=voip \
POSTGRES_PASSWORD=voip_password \
POSTGRES_DB=voip_cloud_pbx \
JWT_ACCESS_SECRET=dev_access_secret_change_me \
JWT_REFRESH_SECRET=dev_refresh_secret_change_me \
JWT_ACCESS_TTL=15m \
JWT_REFRESH_TTL=720h \
go run ./...
```

Then point the gateway to `http://localhost:8081` for auth.

---

## 11. Where to go next

Once everything is running and you can:

- Register/login,
- See Dashboard,
- Send SIP REGISTER/INVITE,
- Observe RTP, QoS, and call analytics,

you can extend:

- **SIP**: real digest auth, INVITE routing to registered Contacts.
- **Media**: call-aware RTP port allocation, recording, WebRTC integration.
- **Dashboard**: active call panels, per-call QoS and recordings, AI summaries.

---



From there you can adjust env variables, ports, or database configuration to fit your own environment.

## Notes and Next Steps

This codebase intentionally focuses on:

- Clear separation of services.
- Simple, testable flows (auth, SIP, RTP, QoS).
- Infrastructure that can be enhanced incrementally.

Possible future enhancements (not yet implemented):

- Full SIP digest authentication.
- Real INVITE routing to registered contacts.
- Persistent QoS and call analytics service.
- Recording and playback.
- WebRTC integration via Pion WebRTC.
