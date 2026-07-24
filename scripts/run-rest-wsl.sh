#!/bin/sh
GO=/usr/local/go/bin/go
cd /mnt/c/Users/fusneif/Personal/cheeky-cart-catalog || exit 1
export PORT=8090
export OTEL_EXPORTER_OTLP_ENDPOINT="${OTEL_EXPORTER_OTLP_ENDPOINT:-127.0.0.1:4317}"
export OTEL_EXPORTER_OTLP_PROTOCOL=grpc
export OTEL_SERVICE_NAME=cheeky-cart-catalog-rest
exec "$GO" run ./cmd/rest
