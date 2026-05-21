# Architecture Overview

VoIP Cloud PBX is organized as a monorepo:

- `services/` contains Go microservices, each as its own Go module.
- `apps/` contains frontend applications built with Next.js.
- `packages/` holds shared Go and TypeScript code.
- `infrastructure/` holds operational configuration for Traefik, databases, caches, monitoring, and Docker orchestration.
- `scripts/` and `docs/` provide tooling and documentation for development and operations.

Detailed diagrams and component descriptions will be added as implementation progresses.
