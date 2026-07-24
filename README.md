# cheeky-cart-catalog (L09 / HW7 checkpoint)

**Branch:** `L09-HW7` — L06/HW5 Catalog (REST + gRPC) plus L09/HW7 telemetry and k6.  
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
| L09 | **HW7** — this branch: structured logs + `/metrics`; **k6** load test (`scripts/catalog-load.js`) |

**Snapshot branches:** [`L06-HW5`](https://github.com/Fusneica-FlorentinCristian/cheeky-cart-catalog/tree/L06-HW5) (REST + gRPC only) · **`L09-HW7`** (this tree) · `main` for later work.

---

## L09 instrumentation + k6 + SigNoz (HW7)

REST server (`cmd/rest`) includes:

- **JSON logs** (`log/slog`) — `request_id`, `method`, `path`, `status`, `duration_ms`
- **Prometheus metrics** at `GET /metrics` — `http_requests_total`, `http_request_duration_seconds`
- **OpenTelemetry traces** (optional) — export to SigNoz when `OTEL_EXPORTER_OTLP_ENDPOINT` is set
- **Correlation ID** — pass `X-Request-ID` header or receive one in the response

Full performance report: [docs/performance-testing.md](docs/performance-testing.md)  
SigNoz setup: [docker/signoz-foundry/README.md](docker/signoz-foundry/README.md)

### Verify instrumentation

```powershell
go run ./cmd/rest
curl http://localhost:8080/health
curl http://localhost:8080/metrics
```

### k6 load test (required for HW7)

Install [k6](https://grafana.com/docs/k6/latest/set-up/install-k6/) once (`winget install GrafanaLabs.k6` on Windows).

With REST running in another terminal:

```powershell
k6 run scripts/catalog-load.js
```

### SigNoz locally (required for HW7)

1. Install [Foundry](https://signoz.io/docs/install/docker/) (`foundryctl`)
2. Sign in to Docker Desktop (org policy may require Autodesk account)
3. Start stack:

```powershell
.\scripts\start-signoz.ps1
```

4. Run REST with trace export (catalog on **8081** while SigNoz UI uses **8080**):

```powershell
$env:PORT = "8081"
$env:OTEL_EXPORTER_OTLP_ENDPOINT = "localhost:4318"
go run ./cmd/rest
curl http://localhost:8081/products
```

5. Open http://localhost:8080 → **Services** → `cheeky-cart-catalog-rest` → **Traces**

### LMS submit (HW7)

1. Push branch **`L09-HW7`** to GitHub
2. Grant collaborator **`popescuag`**
3. LMS comment: repo URL + note HW7 adds telemetry, k6 report, and SigNoz docs
4. Attach **L09-performance-test-report.docx** (from course builder) or link `docs/performance-testing.md`

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
internal/telemetry/            # L09 — logs, metrics, OTel traces
scripts/catalog-load.js        # L09 HW7 — k6 load test
scripts/start-signoz.ps1       # L09 HW7 — start local SigNoz (Foundry)
docker/signoz-foundry/         # L09 HW7 — Foundry casting.yaml + README
docs/performance-testing.md    # L09 HW7 — filled performance report
cmd/rest/main.go               # HTTP JSON — HW5 task 01, HW7 telemetry
cmd/grpc/main.go               # gRPC server — HW5 task 02
```

Shared `internal/catalog` keeps REST and gRPC serving the same data — one service, two transports.
