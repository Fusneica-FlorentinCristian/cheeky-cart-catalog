# Generate gRPC stubs from api/catalog/v1/catalog.proto
# Requires: protoc on PATH, go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
#           go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

$ErrorActionPreference = "Stop"
$root = $PSScriptRoot
$env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" +
            [System.Environment]::GetEnvironmentVariable("Path", "User")

New-Item -ItemType Directory -Force -Path "$root\gen\catalog\v1" | Out-Null

protoc `
  --proto_path="$root\api" `
  --go_out="$root\gen" --go_opt=paths=source_relative `
  --go-grpc_out="$root\gen" --go-grpc_opt=paths=source_relative `
  catalog/v1/catalog.proto

Write-Host "Generated gen/catalog/v1/*.go"
