package telemetry_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Fusneica-FlorentinCristian/cheeky-cart-catalog/internal/telemetry"
)

func TestLoggingMetrics_setsRequestIDAndStatus(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Request-ID") != "client-id" {
			t.Fatalf("expected incoming X-Request-ID to reach handler")
		}
		w.WriteHeader(http.StatusTeapot)
		_, _ = io.WriteString(w, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	req.Header.Set("X-Request-ID", "client-id")
	rec := httptest.NewRecorder()

	telemetry.LoggingMetrics(inner).ServeHTTP(rec, req)

	if rec.Code != http.StatusTeapot {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusTeapot)
	}
	if got := rec.Header().Get("X-Request-ID"); got != "client-id" {
		t.Fatalf("response X-Request-ID = %q, want client-id", got)
	}
}

func TestLoggingMetrics_generatesRequestIDWhenMissing(t *testing.T) {
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()

	telemetry.LoggingMetrics(inner).ServeHTTP(rec, req)

	id := rec.Header().Get("X-Request-ID")
	if id == "" || id == "unknown" {
		t.Fatalf("expected generated request id, got %q", id)
	}
	if len(id) < 8 {
		t.Fatalf("request id too short: %q", id)
	}
}

func TestMetricsHandler_returnsPrometheusText(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()

	telemetry.MetricsHandler().ServeHTTP(rec, req)

	body := rec.Body.String()
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d", rec.Code)
	}
	if !strings.Contains(body, "http_requests_total") {
		t.Fatalf("metrics body missing http_requests_total")
	}
}
