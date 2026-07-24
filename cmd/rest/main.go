// REST API for Catalog bounded context (L06 HW5 task 01, L09 HW7 telemetry).
package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/Fusneica-FlorentinCristian/cheeky-cart-catalog/internal/catalog"
	"github.com/Fusneica-FlorentinCristian/cheeky-cart-catalog/internal/telemetry"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	ctx := context.Background()
	shutdown, err := telemetry.InitTracing(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer telemetry.FlushTracing(ctx, shutdown)

	store := catalog.NewStore()
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, map[string]string{"status": "ok"})
	})
	mux.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		writeJSON(w, store.List())
	})
	mux.HandleFunc("/products/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		id := strings.TrimPrefix(r.URL.Path, "/products/")
		if p, ok := store.Get(id); ok {
			writeJSON(w, p)
			return
		}
		http.NotFound(w, r)
	})
	mux.Handle("/metrics", telemetry.MetricsHandler())

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + strings.TrimPrefix(p, ":")
	}
	otelNote := ""
	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") != "" {
		otelNote = ", OTLP traces"
	}
	log.Printf("catalog REST listening on %s (JSON logs + /metrics%s)", addr, otelNote)
	handler := telemetry.TraceMiddleware(telemetry.LoggingMetrics(mux))
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatal(err)
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
