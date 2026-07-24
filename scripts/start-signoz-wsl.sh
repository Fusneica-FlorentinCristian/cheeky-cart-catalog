#!/usr/bin/env bash
# Start SigNoz via native Docker Engine in WSL (bypasses Docker Desktop org policy).
set -euo pipefail

export PATH="${HOME}/.local/bin:${PATH}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/.." && pwd)"
FOUNDRY_DIR="${REPO_ROOT}/docker/signoz-foundry"

if ! command -v foundryctl >/dev/null 2>&1; then
  echo "Installing Foundry..."
  curl -fsSL https://signoz.io/foundry.sh | bash
  export PATH="${HOME}/.local/bin:${PATH}"
fi

cd "${FOUNDRY_DIR}"
if [[ ! -f pours/deployment/compose.yaml ]]; then
  foundryctl forge -f casting.yaml
fi

cd pours/deployment
docker compose up -d

echo ""
echo "SigNoz UI: http://localhost:8080"
echo "OTLP HTTP: localhost:4318"
echo ""
echo "From Windows PowerShell (catalog on 8081):"
echo '  $env:PORT="8081"; $env:OTEL_EXPORTER_OTLP_ENDPOINT="localhost:4318"; go run ./cmd/rest'
