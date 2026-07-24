# cheeky-cart-catalog (HW5 checkpoint)

**Branch:** `hw5` — REST + gRPC only (no later homework additions).  
**Module:** `github.com/Fusneica-FlorentinCristian/cheeky-cart-catalog`

Cheeky Cart **Catalog** service — **FR-01** browse/list, **FR-02** product detail.

| Task | Transport | Run | Port |
|------|-----------|-----|------|
| 01 | REST | `go run ./cmd/rest` | `:8080` |
| 02 | gRPC | `go run ./cmd/grpc` | `:50051` |

## REST

```powershell
go run ./cmd/rest
curl http://localhost:8080/health
curl http://localhost:8080/products
curl http://localhost:8080/products/1
```

## gRPC

```powershell
go run ./cmd/grpc
grpcurl -plaintext localhost:50051 catalog.v1.CatalogService/ListProducts
grpcurl -plaintext -d "{\"id\":\"1\"}" localhost:50051 catalog.v1.CatalogService/GetProduct
```

## Layout

```
api/catalog/v1/catalog.proto
cmd/rest/main.go
cmd/grpc/main.go
gen/catalog/v1/
internal/catalog/store.go
```

Later work (telemetry, k6) lives on `main`, not this branch.
