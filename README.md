# community-helper-app
Android-first platform that pairs help seekers with nearby helpers for urgent or scheduled support, blending Flutter/Firebase, Twilio SMS, and Stripe escrow to deliver trusted community assistance with real-time matching, chat, and impact tracking.

## Backend Service Skeleton
- Gin-powered HTTP API under `cmd/server` aligned with the `docs/API.md` surface.
- Request/response models live in `internal/domain/models`; in-memory service implementations under `internal/domain/services/memory` simulate Firebase/Stripe/Twilio behaviour.
- Handlers in `internal/api/handlers` include validation, auth guard, and business logic pipelines.
- Start the server with:
  ```bash
  HTTP_PORT=8080 go run ./cmd/server
  ```
- Configure a writable Go build cache if required by your environment:
  ```bash
  export GOCACHE=$(pwd)/.cache
  export GOMODCACHE=$(pwd)/.modcache
  ```
- Run unit tests (also used in CI):
  ```bash
  GOCACHE=$(pwd)/.cache GOMODCACHE=$(pwd)/.modcache go test ./...
  ```
