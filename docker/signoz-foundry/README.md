# SigNoz local stack (L09 HW7)

Observability backend for **OpenTelemetry traces** exported by `cmd/rest`.

## Prerequisites

- Docker Engine 20.10+ with Compose v2 plugin
- [Foundry](https://github.com/SigNoz/foundry) (`foundryctl`)
- At least 4 GB RAM allocated to Docker

Install Foundry (Linux/macOS/WSL):

```bash
curl -fsSL https://signoz.io/foundry.sh | bash
```

On Windows PowerShell, see [foundry getting-started](https://github.com/SigNoz/foundry/blob/main/docs/getting-started.md).

## Start SigNoz

### Option A — WSL native Docker (recommended on Windows)

Bypasses Docker Desktop Autodesk org login. Requires Ubuntu WSL with **native Docker Engine** (you already have this if `wsl -d Ubuntu-22.04 docker version` shows a Linux server).

```powershell
.\scripts\start-signoz-wsl.ps1
```

Run the catalog **inside WSL** for reliable OTLP export:

```powershell
wsl -d Ubuntu-22.04 sh /mnt/c/Users/fusneif/Personal/cheeky-cart-catalog/scripts/run-rest-wsl.sh
```

Catalog listens on **8090** (SigNoz UI keeps **8080**).

### Option B — Docker Desktop (if org login works)

```powershell
.\scripts\start-signoz.ps1
```

Or manually:

```powershell
cd docker/signoz-foundry
foundryctl gauge -f casting.yaml
foundryctl forge -f casting.yaml
cd pours/deployment
docker compose up -d
```

- **SigNoz UI:** http://localhost:8080
- **OTLP gRPC:** `localhost:4317` (prefer `127.0.0.1:4317` from Windows)
- **OTLP HTTP:** `localhost:4318`

First launch may take 1–2 minutes while ClickHouse becomes healthy.

### First-run org / OTLP gotcha

If **no traces** appear and OTLP `:4317`/`:4318` refuse connections, the ingester is likely still on **nop** pipelines because OpAMP has no org:

```text
cannot create agent without orgId
```

**Fix:** enable the root user in `compose.override.yaml` (already in this repo) and restart the stack:

```yaml
SIGNOZ_USER_ROOT_ENABLED: "true"
SIGNOZ_USER_ROOT_EMAIL: admin@cheeky-cart.local
SIGNOZ_USER_ROOT_PASSWORD: CheekyCart-HW7-Local
SIGNOZ_USER_ROOT_ORG_NAME: default
```

Copy the override into `pours/deployment/` after `foundryctl forge`, then `docker compose up -d`. After `setupCompleted: true`, restart the ingester:

```bash
docker restart signoz-ingester-1
docker exec signoz-ingester-1 grep otlp /var/tmp/collector-config.yaml   # should list otlp receivers, not nop
```

Alternative: complete signup once at http://localhost:8080/signup.

## Run catalog with traces

**WSL (recommended):**

```powershell
wsl -d Ubuntu-22.04 sh scripts/run-rest-wsl.sh
curl http://localhost:8090/products
```

**Windows PowerShell** (use gRPC + `127.0.0.1`; HTTP/localhost can fail across WSL port forwarding):

```powershell
$env:PORT = "8081"
$env:OTEL_EXPORTER_OTLP_ENDPOINT = "127.0.0.1:4317"
$env:OTEL_EXPORTER_OTLP_PROTOCOL = "grpc"
$env:OTEL_SERVICE_NAME = "cheeky-cart-catalog-rest"
go run ./cmd/rest
```

In SigNoz UI → **Services** → `cheeky-cart-catalog-rest` → **Traces** to inspect spans.

## Stop

```powershell
cd docker/signoz-foundry/pours/deployment
docker compose down
```

## Windows note

**Do not use Docker Desktop** if your org blocks image pulls — use **WSL native Docker** (`.\scripts\start-signoz-wsl.ps1`) instead.

SigNoz recommends native Docker Engine inside WSL 2. If ClickHouse Keeper crashes (exit 139) under Docker Desktop, use the WSL path. See [SigNoz Docker install](https://signoz.io/docs/install/docker/).
