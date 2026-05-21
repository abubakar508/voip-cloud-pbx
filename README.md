# VoIP Cloud PBX – Monorepo

A multi-service VoIP Cloud PBX platform built with Go (backend microservices) and Next.js (frontends). It includes:

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

## Prerequisites

- Docker and Docker Compose
- Go 1.21+ (for local development)
- Node.js 18+ and pnpm/npm (for frontend development)

---

## Environment

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

## Running with Docker Compose

From the repo root:

```bash
# Build images
docker compose build

# Start core infra and services
docker compose up -d \
  postgres redis nats traefik \
  auth-service api-gateway websocket-service \
  sip-service media-service \
  dashboard-web
```

Check status:

```bash
docker compose ps
```

Tail logs:

```bash
docker compose logs -f auth-service api-gateway websocket-service sip-service media-service
```

---

## Basic Workflows

### 1. Register and Login (Auth + Dashboard)

1. Register a tenant admin:

   ```bash
   curl -X POST http://localhost:8080/auth/register \
     -H "Content-Type: application/json" \
     -d '{
       "tenantName": "Demo Tenant",
       "email": "admin@demo.local",
       "password": "Passw0rd!",
       "displayName": "Demo Admin"
     }'
   ```

2. Login:

   ```bash
   curl -X POST http://localhost:8080/auth/login \
     -H "Content-Type: application/json" \
     -d '{
       "email": "admin@demo.local",
       "password": "Passw0rd!"
     }'
   ```

3. Open the dashboard in a browser and log in using the same credentials.

### 2. SIP Registration and Calls

1. Create SIP accounts in Postgres (e.g. using `psql` or migrations) with bcrypt-hashed passwords.

2. Configure a SIP client (e.g., Linphone, Zoiper):

   - SIP server: your host IP
   - Port: 5060
   - Transport: UDP
   - Username: your SIP account username
   - Password: matching password
   - Realm: as needed

3. REGISTER:

   - sip-service listens on UDP 5060.
   - REGISTER will be parsed and validated; successful registrations are stored in Redis.

4. INVITE/BYE/CANCEL:

   - When calls are made, sip-service creates call sessions and publishes events to NATS.
   - media-service subscribes and tracks calls.

### 3. RTP + QoS

1. Send RTP to media-service (UDP 40000) from test tools or softphones.

2. Visit:

   - `http://localhost:8082/calls` – see tracked calls.
   - `http://localhost:8082/qos` – see per-stream QoS.

3. To test RTP forwarding:

   - Configure endpoints for a call via:

     ```bash
     curl -X POST http://localhost:8082/calls/<callId>/endpoints \
       -H "Content-Type: application/json" \
       -d '{"aAddr":"127.0.0.1:5000","bAddr":"127.0.0.1:5002"}'
     ```

   - Send RTP to media-service from port 5000 and listen on port 5002.

### 4. Dashboard: Calls and QoS

On the dashboard home page (after login), you can see:

- A table of active/ended calls fetched from media-service `/calls`.
- A table of QoS streams fetched from `/qos`.

---

## Development

### Go Services (local)

You can run a single service locally (outside Docker) if you have Postgres/Redis/NATS running:

```bash
cd services/auth-service
APP_ENV=development \
POSTGRES_HOST=localhost POSTGRES_PORT=5432 \
POSTGRES_USER=voip POSTGRES_PASSWORD=voip_password POSTGRES_DB=voip_cloud_pbx \
JWT_ACCESS_SECRET=dev_access JWT_REFRESH_SECRET=dev_refresh JWT_ACCESS_TTL=15m JWT_REFRESH_TTL=720h \
go run ./...
```

Do similarly for other services, adjusting env vars as needed.

### Frontend

```bash
cd apps/dashboard-web
npm install
npm run dev
```

Then open `http://localhost:3000`.

Make sure your API and media endpoints match the URLs the app expects (either via `.env.local` or hardcoded during development).

---

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
