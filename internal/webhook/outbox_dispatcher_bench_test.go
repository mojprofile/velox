package webhook_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/sagarsuperuser/velox/internal/platform/postgres"
	"github.com/sagarsuperuser/velox/internal/webhook"
)

// BenchmarkOutboxDispatcher_BeginTxError isolates the dispatcher's failure
// path from database latency. A cancelled context makes BeginTx fail before
// the driver opens a connection; the same benchmark runs before and after
// adding the enqueue-failure metric.
func BenchmarkOutboxDispatcher_BeginTxError(b *testing.B) {
	pool, err := sql.Open("pgx", "postgres://unused:unused@127.0.0.1:1/unused?sslmode=disable")
	if err != nil {
		b.Fatal(err)
	}
	b.Cleanup(func() { _ = pool.Close() })
	dispatcher := webhook.NewOutboxDispatcher(webhook.NewOutboxStore(postgres.NewDB(pool, 0)))
	ctx, cancel := context.WithCancel(postgres.WithLivemode(context.Background(), false))
	cancel()
	const eventType = "invoice.finalized"
	payload := map[string]any{"invoice_id": "inv_benchmark"}
	if err := dispatcher.Dispatch(ctx, "tenant_benchmark", eventType, payload); !errors.Is(err, context.Canceled) {
		b.Fatalf("dispatch with cancelled context: got %v, want context.Canceled", err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := dispatcher.Dispatch(ctx, "tenant_benchmark", eventType, payload); err == nil {
			b.Fatal("expected transaction start error")
		}
	}
}
