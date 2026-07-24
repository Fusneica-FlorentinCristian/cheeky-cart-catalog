# Performance testing results and recommendations

Instructor template: [course/homework/L09/templates/hw9-performance-testing.docx](https://github.com/Fusneica-FlorentinCristian/software-architecture/blob/master/course/homework/L09/templates/hw9-performance-testing.docx)

Full report (filled): [docs/performance-testing.md](performance-testing.md)

## Quick results (k6 Run 1)

| Metric | Value |
|--------|-------|
| Throughput | 811 req/s (40,647 requests) |
| p95 latency | 1.07 ms |
| Errors | 0% |
| NFR-02 (p95 &lt; 2000 ms) | Pass |

```powershell
go run ./cmd/rest
k6 run scripts/catalog-load.js
```

## SigNoz (traces)

See [docker/signoz-foundry/README.md](../docker/signoz-foundry/README.md).

**Verified OTLP export:** [docs/signoz-otlp-verification.md](signoz-otlp-verification.md) (ClickHouse query snapshots).

```powershell
.\scripts\start-signoz-wsl.ps1
wsl -d Ubuntu-22.04 sh scripts/run-rest-wsl.sh
curl http://localhost:8090/products
```

Open http://localhost:8080 → **Services** → `cheeky-cart-catalog-rest`.
