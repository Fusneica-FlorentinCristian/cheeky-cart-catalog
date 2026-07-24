# Start SigNoz via Foundry (L09 HW7).
# Usage: .\scripts\start-signoz.ps1
$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $PSScriptRoot
$Dir = Join-Path $Root "docker\signoz-foundry"
$Casting = Join-Path $Dir "casting.yaml"

if (-not (Get-Command foundryctl -ErrorAction SilentlyContinue)) {
    Write-Error "foundryctl not found. Install Foundry: https://signoz.io/docs/install/docker/"
}

Push-Location $Dir
try {
    foundryctl gauge -f $Casting
    foundryctl forge -f $Casting
    Push-Location "pours\deployment"
    docker compose up -d
    Write-Host ""
    Write-Host "SigNoz UI: http://localhost:8080"
    Write-Host "OTLP HTTP: localhost:4318  |  OTLP gRPC: localhost:4317"
    Write-Host ""
    Write-Host 'With SigNoz running, start REST on 8081: $env:PORT="8081"; $env:OTEL_EXPORTER_OTLP_ENDPOINT="localhost:4318"; go run ./cmd/rest'
}
finally {
    Pop-Location
    Pop-Location
}
