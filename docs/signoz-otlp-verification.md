# SigNoz OTLP verification (L09 HW7)

Captured **2026-07-30** after fixing the ingester **nop-pipeline** issue (root user org bootstrap in `compose.override.yaml`).

Screenshots taken at **1920×1080** after a **k6 load run with OTLP export enabled** (~79k spans). UI charts stay empty if you only send a handful of manual requests.

## Environment

| Item | Value |
|------|-------|
| Stack | SigNoz Foundry via WSL native Docker |
| SigNoz UI | http://localhost:8080 |
| Catalog REST | Windows `:8090`, `OTEL_EXPORTER_OTLP_ENDPOINT=127.0.0.1:4317`, `OTEL_EXPORTER_OTLP_PROTOCOL=grpc` |
| Service name | `cheeky-cart-catalog-rest` |

**Prerequisite:** copy `compose.override.yaml` into `docker/signoz-foundry/pours/deployment/` before `docker compose up -d`, then `docker restart signoz-ingester-1`.

## SigNoz UI — Services list

`cheeky-cart-catalog-rest` appears under **Services** with latency and ops/s derived from exported spans:

![SigNoz Services — cheeky-cart-catalog-rest](images/signoz-services.png)

## SigNoz UI — Service metrics & key operations

Drill into the service to see latency/rate charts and per-endpoint span counts (`GET /products`, `GET /health`, etc.):

![SigNoz service overview — latency, rate, Apdex, key operations](images/signoz-service-detail.png)

Close-up of the **Key Operations** table (span names = HTTP method + path from `otelhttp` middleware):

![SigNoz key operations — GET /products, GET /health, GET /products/1](images/signoz-service-operations.png)

## ClickHouse — backend proof

Same data visible in ClickHouse (`signoz_traces.signoz_index_v3`):

![ClickHouse trace count and span breakdown](images/clickhouse-traces-proof.png)

<details>
<summary>Raw ClickHouse queries (text)</summary>

```sql
SELECT count() AS total_traces FROM signoz_traces.signoz_index_v3
-- 79681

SELECT serviceName, name, count() AS spans
FROM signoz_traces.signoz_index_v3
GROUP BY serviceName, name
ORDER BY spans DESC
```

```
   ┌─serviceName──────────────┬─name────────────┬─spans─┐
1. │ cheeky-cart-catalog-rest │ GET /health     │ 26563 │
2. │ cheeky-cart-catalog-rest │ GET /products   │ 26563 │
3. │ cheeky-cart-catalog-rest │ GET /products/1 │ 26555 │
   └──────────────────────────┴─────────────────┴───────┘
```

</details>

## Ingester config (OTLP active)

After org bootstrap + ingester restart, runtime config uses real `otlp` receivers (not `nop`):

```yaml
receivers:
    otlp:
        protocols:
            grpc:
            ...
pipelines:
    traces:
        receivers:
            - otlp
```

## Reproduce

```powershell
cd docker/signoz-foundry/pours/deployment
docker compose up -d
# ensure compose.override.yaml (SIGNOZ_USER_ROOT_*) is present
docker restart signoz-ingester-1

cd ../../..
$env:PORT="8090"
$env:OTEL_EXPORTER_OTLP_ENDPOINT="127.0.0.1:4317"
$env:OTEL_EXPORTER_OTLP_PROTOCOL="grpc"
go run ./cmd/rest

# In another terminal — generate enough spans for UI charts
k6 run -e BASE_URL=http://localhost:8090 scripts/catalog-load.js
```

Open http://localhost:8080 → **Services** → `cheeky-cart-catalog-rest`.
