# Start SigNoz via WSL native Docker (bypasses Docker Desktop org login).
# Prereq: Ubuntu WSL with native Docker Engine (not Docker Desktop integration).
# Usage: .\scripts\start-signoz-wsl.ps1

$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $PSScriptRoot
$Foundry = Join-Path $Root "docker\signoz-foundry"
$Foundryctl = "$env:LOCALAPPDATA\foundry\foundry_windows_amd64\bin\foundryctl.exe"

if (-not (Test-Path $Foundryctl)) {
    Write-Host "Installing Foundry in WSL..."
    wsl -d Ubuntu-22.04 bash -c "curl -fsSL https://signoz.io/foundry.sh | bash"
    $Foundryctl = "/root/.local/bin/foundryctl"
} else {
    $Foundryctl = "/root/.local/bin/foundryctl"
    wsl -d Ubuntu-22.04 bash -c "test -x $Foundryctl || (curl -fsSL https://signoz.io/foundry.sh | bash)"
}

$WslFoundry = "/mnt/c/Users/fusneif/Personal/cheeky-cart-catalog/docker/signoz-foundry"
wsl -d Ubuntu-22.04 bash -c "export PATH=/root/.local/bin:`$PATH; cd '$WslFoundry' && foundryctl forge -f casting.yaml && cp compose.override.yaml pours/deployment/compose.override.yaml && cd pours/deployment && docker compose -f compose.yaml -f compose.override.yaml up -d"

Write-Host ""
Write-Host "SigNoz UI: http://localhost:8080"
Write-Host "OTLP gRPC: localhost:4317  |  OTLP HTTP: localhost:4318"
Write-Host ""
Write-Host "Run catalog with traces (WSL — recommended for OTLP):"
Write-Host "  wsl -d Ubuntu-22.04 sh /mnt/c/Users/fusneif/Personal/cheeky-cart-catalog/scripts/run-rest-wsl.sh"
Write-Host ""
Write-Host "Or from Windows (gRPC, use 127.0.0.1 not localhost):"
Write-Host '  $env:PORT="8081"; $env:OTEL_EXPORTER_OTLP_ENDPOINT="127.0.0.1:4317"; $env:OTEL_EXPORTER_OTLP_PROTOCOL="grpc"; go run ./cmd/rest'
