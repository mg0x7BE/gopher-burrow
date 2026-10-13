# Inventory service

A **small** service with an _in-memory_ store and `go test` coverage.

## Setup

1. Install [Go](https://go.dev/dl/) 1.27 or newer.
2. Run the build:

   ```bash
   go run ./cmd/build
   ```

> Tip: set `LOG_LEVEL=debug` to see every request.

- [x] Store
- [ ] Persistence

| Endpoint | Method | Auth |
|----------|--------|------|
| `/items` | GET    | no   |
| `/items` | POST   | yes  |

---

![Gopher](https://go.dev/images/gophers/ladder.svg)
