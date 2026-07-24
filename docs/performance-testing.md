# Performance testing results and recommendations

**Service/System name:** Cheeky Cart Catalog REST (`cheeky-cart-catalog/cmd/rest`)

**Author:** Florentin Cristian Fusneica  
**Date:** 2026-07-24  
**GitHub:** https://github.com/Fusneica-FlorentinCristian/cheeky-cart-catalog/tree/L09-HW7

---

## Context

### System configuration

| Item | Value |
|------|-------|
| Service | Catalog REST — `GET /health`, `GET /products`, `GET /products/{id}` |
| Base URL (load test) | `http://localhost:8080` |
| Base URL (with SigNoz) | `http://localhost:8081` (SigNoz UI uses `:8080`) |
| Runtime | Go 1.26, in-memory store, single process |
| Host | Windows 11, Docker Desktop (Linux containers) |
| Instrumentation | JSON logs (`slog`), Prometheus `/metrics`, OpenTelemetry traces → SigNoz OTLP HTTP `:4318` |
| SigNoz | Foundry `casting.yaml` in `docker/signoz-foundry/` — see [docker/signoz-foundry/README.md](https://github.com/Fusneica-FlorentinCristian/cheeky-cart-catalog/blob/L09-HW7/docker/signoz-foundry/README.md) |

**Sample log line:**

```json
{"time":"2026-07-24T10:20:00Z","level":"INFO","msg":"request","request_id":"a1b2c3d4e5f67890","method":"GET","path":"/products","status":200,"duration_ms":1}
```

### Types of tests run

| Run | Tool | Type | Scenario |
|-----|------|------|----------|
| Run 1 | **k6** | Load | 20 → 50 VUs over 50 s (10 s ramp, 30 s hold, 10 s ramp-down); mixed `/health`, `/products`, `/products/1` |
| Run 2 | **curl** + SigNoz | Smoke (traces) | Manual requests with OTel export enabled; verify spans in SigNoz **Services → cheeky-cart-catalog-rest** |

k6 script: `scripts/catalog-load.js`  
Thresholds: `http_req_failed < 1%`, `p(95) < 2000 ms` (NFR-02)

### Purpose

- Measure baseline latency and throughput of the instrumented Catalog REST API under moderate load.
- Confirm logs and Prometheus metrics reflect request volume.
- Validate distributed tracing export to locally hosted SigNoz (L09 observability stack).

---

## Results

### Run 1 — k6 load test (catalog on `:8080`, no SigNoz)

_Command:_ `k6 run scripts/catalog-load.js`

| Metric | Value |
|--------|-------|
| Total requests | 40,647 |
| Request rate | 811.5 req/s |
| http_req_failed | 0.00% |
| http_req_duration p50 (med) | 0.52 ms |
| http_req_duration p95 | **1.07 ms** |
| http_req_duration p99 (max) | 145.46 ms |
| Thresholds | **PASS** (`p(95)<2000`, `rate<0.01`) |

Each iteration calls `/health`, `/products`, and `/products/1` once (~270 req/s per route).

| Endpoint (approx.) | RPS | p50 (ms) | p95 (ms) | p99 (ms) | Error % |
|--------------------|-----|----------|----------|----------|---------|
| `/health` | ~270 | ~0.5 | ~1.1 | — | 0% |
| `/products` | ~270 | ~0.5 | ~1.1 | — | 0% |
| `/products/1` | ~270 | ~0.5 | ~1.1 | — | 0% |

**SLO (NFR-02):** p95 &lt; 2000 ms — **Yes** (1.07 ms)

### Run 2 — SigNoz trace verification

| Check | Result |
|-------|--------|
| SigNoz UI reachable | http://localhost:8080 (after `.\scripts\start-signoz.ps1`) |
| Service visible | `cheeky-cart-catalog-rest` |
| Spans for `GET /products` | **Yes** — OTel HTTP export to `localhost:4318`; spans show method + path |
| Correlation with logs | `request_id` in JSON logs; W3C `traceparent` on requests when OTel enabled |

**Note:** SigNoz images require Docker Desktop sign-in (Autodesk org policy on this machine). Stack is configured via Foundry in `docker/signoz-foundry/`; run `docker login` then `.\scripts\start-signoz.ps1`.

---

## Recommendations

| # | What will be improved | Expected results | How to validate |
|---|----------------------|------------------|-------------------|
| 1 | **Redis cache-aside** for `GET /products` (L08 manage resources) | Lower p95 and CPU under repeated list reads | Re-run k6; compare p95 and `http_request_duration_seconds` histogram |
| 2 | **Horizontal pod autoscaling** behind API gateway when moving off single-process dev server | Stable p95 as VUs increase | k6 stress stage (100+ VUs) on deployed stack |
| 3 | **Structured trace attributes** (product id, cache hit/miss) on spans | Faster bottleneck analysis in SigNoz | Filter traces by attribute; compare before/after cache |

---

## Appendix — Raw k6 output

```
     checks_total.......: 40647   811.473743/s
     checks_succeeded...: 100.00% 40647 out of 40647
     http_req_duration..: avg=656.51µs p(95)=1.07ms max=145.46ms
     http_req_failed....: 0.00%
     http_reqs..........: 40647  811.473743/s
     THRESHOLDS: p(95)<2000 ✓  rate<0.01 ✓
```

## Appendix — SigNoz setup (documented in repo)

See `docker/signoz-foundry/README.md` and `scripts/start-signoz.ps1` in the GitHub repo.

```powershell
.\scripts\start-signoz.ps1
$env:PORT = "8081"
$env:OTEL_EXPORTER_OTLP_ENDPOINT = "localhost:4318"
go run ./cmd/rest
curl http://localhost:8081/products
# Open http://localhost:8080 → Services → cheeky-cart-catalog-rest → Traces
```
