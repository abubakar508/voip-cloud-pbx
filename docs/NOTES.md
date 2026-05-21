# Developer Notes

- Use Go 1.26 for all backend services.
- Use `docker compose` (not `docker-compose`) when running the stack.
- Each service should be independently buildable and testable.
- Shared functionality goes into `packages/shared-go` and `packages/shared-ts` to avoid duplication.
