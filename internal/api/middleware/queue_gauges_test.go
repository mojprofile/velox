package middleware_test

import (
	"errors"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	mw "github.com/sagarsuperuser/velox/internal/api/middleware"
)

func TestQueueDepthGauges_QueryFailure(t *testing.T) {
	registry := prometheus.NewRegistry()
	previousRegisterer := prometheus.DefaultRegisterer
	prometheus.DefaultRegisterer = registry
	t.Cleanup(func() {
		prometheus.DefaultRegisterer = previousRegisterer
	})

	queryErr := errors.New("database unavailable")
	mw.RegisterQueueDepthGauges(time.Minute, func(string) (float64, error) {
		return 0, queryErr
	})

	families, err := registry.Gather()
	if err != nil {
		t.Fatalf("gather metrics: %v", err)
	}
	if len(families) == 0 {
		t.Fatal("no queue gauges registered")
	}
	for _, family := range families {
		for _, metric := range family.GetMetric() {
			gauge := metric.GetGauge()
			if gauge == nil {
				t.Errorf("%s: expected a gauge", family.GetName())
				continue
			}
			if got := gauge.GetValue(); got != -1 {
				t.Errorf("%s: got %v, want -1 on query failure", family.GetName(), got)
			}
		}
	}
}
