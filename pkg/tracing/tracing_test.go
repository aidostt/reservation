package tracing

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"go.opentelemetry.io/otel"
)

// TestInit_ExportsSpans verifies the OTLP export path end to end: it initialises
// tracing against a real collector, emits a span, and confirms the collector
// received it. It is skipped unless both OTEL_EXPORTER_OTLP_ENDPOINT and
// JAEGER_QUERY_URL are set, so it stays inert in CI when no collector is around.
func TestInit_ExportsSpans(t *testing.T) {
	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") == "" || os.Getenv("JAEGER_QUERY_URL") == "" {
		t.Skip("set OTEL_EXPORTER_OTLP_ENDPOINT and JAEGER_QUERY_URL to run")
	}
	ctx := context.Background()
	serviceName := fmt.Sprintf("tracing-verify-%d", time.Now().UnixNano())

	shutdown, err := Init(ctx, serviceName)
	if err != nil {
		t.Fatalf("Init: %v", err)
	}
	_, span := otel.Tracer("verify").Start(ctx, "verify-span")
	span.End()
	if err := shutdown(ctx); err != nil {
		t.Fatalf("shutdown/flush: %v", err)
	}

	queryURL := os.Getenv("JAEGER_QUERY_URL")
	deadline := time.Now().Add(20 * time.Second)
	for {
		if collectorHasService(queryURL, serviceName) {
			return // the span reached the collector
		}
		if time.Now().After(deadline) {
			t.Fatalf("service %q never appeared in the collector", serviceName)
		}
		time.Sleep(time.Second)
	}
}

func collectorHasService(queryURL, serviceName string) bool {
	resp, err := http.Get(queryURL + "/api/services")
	if err != nil {
		return false
	}
	defer func() { _ = resp.Body.Close() }()
	var body struct {
		Data []string `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return false
	}
	for _, s := range body.Data {
		if s == serviceName {
			return true
		}
	}
	return false
}
