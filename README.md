# cheeky-cart-catalog

**Module:** `github.com/Fusneica-FlorentinCristian/cheeky-cart-catalog`

Cheeky Cart **Catalog** microservice (L03 bounded context). Implements **FR-01** (browse/list) and **FR-02** (product detail).

| HW5 task | Transport | Entry point | Port |
|----------|-----------|-------------|------|
| 01 | REST (JSON) | `cmd/rest` | `:8080` |
| 02 | gRPC (protobuf) | `cmd/grpc` | `:50051` |

This folder is **self-contained** — you can run it inside the course repo today, then copy it to your own GitHub repo for LMS submit.

---

## Use in current form (course repo)

From the **software-architecture** repo root:

```powershell
cd course/homework/L06/cheeky-cart-catalog
```

**Terminal 1 — REST**

```powershell
go run ./cmd/rest
```

```powershell
curl http://localhost:8080/health
curl http://localhost:8080/products
curl http://localhost:8080/products/1
curl http://localhost:8080/metrics
```

**Terminal 2 — gRPC** (separate terminal)

```powershell
cd course/homework/L06/cheeky-cart-catalog
go run ./cmd/grpc
```

Test with [grpcurl](https://github.com/fullstorydev/grpcurl) (install once if needed):

```powershell
grpcurl -plaintext localhost:50051 catalog.v1.CatalogService/ListProducts
grpcurl -plaintext -d "{\"id\":\"1\"}" localhost:50051 catalog.v1.CatalogService/GetProduct
```

No database, no `.env`, no GitHub — `go mod tidy` runs automatically on first `go run`.

---

## What to extract to your GitHub repo

Copy the **entire** `cheeky-cart-catalog/` directory into a new repo (repo root = this folder contents, or repo named `cheeky-cart-catalog` with this tree at root).

### Include (all of these)

```
cheeky-cart-catalog/
  api/catalog/v1/catalog.proto    # gRPC contract
  gen/catalog/v1/*.go             # generated stubs (already committed — clone-and-run)
  internal/catalog/store.go       # shared in-memory catalog
  internal/telemetry/middleware.go # L09 — logs + Prometheus metrics
  scripts/catalog-load.js          # L09 HW7 — k6 load test
  cmd/rest/main.go                # REST server
  cmd/grpc/main.go                # gRPC server
  go.mod
  go.sum
  README.md                       # this file (edit module path note after copy)
  gen.ps1                         # optional: regenerate stubs after .proto edits
```

### Do not copy from the course repo

| Path | Why |
|------|-----|
| `../starter-rest/` | Legacy REST-only stub — superseded |
| `course/homework/completed/L06/*.docx` | Reference doc only — not part of the Go repo |
| `.venv/`, `tests/` | Course tooling |

### After copy — rename module path

1. Choose GitHub path, e.g. `github.com/yourname/cheeky-cart-catalog`
2. Replace `YOUR_USER` in:
   - `go.mod` (`module` line)
   - `api/catalog/v1/catalog.proto` (`option go_package`)
   - imports in `cmd/rest/main.go`, `cmd/grpc/main.go`
3. If you change `go_package`, rerun `gen.ps1` (needs `protoc` on PATH) or keep `gen/` as-is if path matches your module.

### Verify extracted repo

```powershell
go build ./...
go test ./...
go run ./cmd/rest
k6 run scripts/catalog-load.js
go run ./cmd/grpc
```

Add a one-line note at the top of this README with your real module path before push.

---

## LMS submit (HW5)

1. Push extracted repo to GitHub (public or private per cohort)
2. Grant collaborator access to **`popescuag`** (instructor)
3. LMS comment: repo URL + one line, e.g. `Catalog service — REST :8080, gRPC :50051`
4. **No ARD PDF** for this homework — code only ([INSTRUCTOR-ARD-MAP.md](../../INSTRUCTOR-ARD-MAP.md))

---

## How this connects to other lessons

| Lesson | Link |
|--------|------|
| L03 | **Catalog** bounded context — service split |
| L04 | REST chosen for browser; gRPC optional for internal calls |
| L07 | Document `GET /v1/products` in ARD **Data Flows/APIs** — same contract as REST here |
| L09 | **HW7** — structured logs + `/metrics`; **k6** load test (`scripts/catalog-load.js`) |

---

## L09 instrumentation + k6 (HW7)

REST server (`cmd/rest`) includes:

- **JSON logs** (`log/slog`) — `request_id`, `method`, `path`, `status`, `duration_ms`
- **Prometheus metrics** at `GET /metrics` — `http_requests_total`, `http_request_duration_seconds`
- **Correlation ID** — pass `X-Request-ID` header or receive one in the response

### Verify instrumentation

```powershell
go run ./cmd/rest
```

```powershell
curl http://localhost:8080/health
curl http://localhost:8080/metrics
```

JSON log lines appear on stdout for each request.

### k6 load test (required for HW7)

Install [k6](https://grafana.com/docs/k6/latest/set-up/install-k6/) once (`choco install k6` on Windows).

With REST running in another terminal:

```powershell
k6 run scripts/catalog-load.js
```

Optional base URL override:

```powershell
$env:BASE_URL="http://localhost:8080"
k6 run scripts/catalog-load.js
```

Copy k6 summary output into your [performance test report](../L09/performance-test-report.md) (course repo) or `docs/performance-test-report.md` in your GitHub repo.

### LMS submit (HW7)

1. Push instrumented repo to GitHub (same repo as HW5 or a new branch)
2. Grant collaborator **`popescuag`**
3. LMS comment: repo URL + note that HW7 adds telemetry + k6 results in report
4. **No ARD PDF** — code + performance report

---

## Prerequisites

- **Go 1.22+**
- **Regenerating protobuf** (only if you edit `.proto`): `protoc`, `protoc-gen-go`, `protoc-gen-go-grpc` — see `gen.ps1`

---

## Project layout

```
api/catalog/v1/catalog.proto   # ListProducts, GetProduct RPCs
gen/catalog/v1/                # generated Go (committed)
internal/catalog/store.go      # in-memory products (v1)
internal/telemetry/            # L09 — request logs + Prometheus metrics
scripts/catalog-load.js        # L09 HW7 — k6 load test
cmd/rest/main.go               # HTTP JSON — HW5 task 01, HW7 telemetry
cmd/grpc/main.go               # gRPC server — HW5 task 02
```

Shared `internal/catalog` keeps REST and gRPC serving the same data — one service, two transports.
