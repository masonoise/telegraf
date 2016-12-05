package prometheus

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/influxdata/telegraf/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const sampleTextFormat = `# HELP go_gc_duration_seconds A summary of the GC invocation durations.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 0.00010425500000000001
go_gc_duration_seconds{quantile="0.25"} 0.000139108
go_gc_duration_seconds{quantile="0.5"} 0.00015749400000000002
go_gc_duration_seconds{quantile="0.75"} 0.000331463
go_gc_duration_seconds{quantile="1"} 0.000667154
go_gc_duration_seconds_sum 0.0018183950000000002
go_gc_duration_seconds_count 7
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 15
# HELP test:classes:total Number of classes that currently exist.
# TYPE test:classes:total gauge
test:classes:total 0.0
`

func TestPrometheusGeneratesMetrics(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, sampleTextFormat)
	}))
	defer ts.Close()

	p := &Prometheus{
		Urls: []string{ts.URL},
	}

	var acc testutil.Accumulator

	err := p.Gather(&acc)
	require.NoError(t, err)

	assert.True(t, acc.HasFloatField("go_gc_duration_seconds", "count"))
	assert.True(t, acc.HasFloatField("go_goroutines", "gauge"))
	assert.True(t, acc.HasFloatField("test:classes:total", "gauge"))
}

// Test that the options to change the separator to a period, and add a prefix
// to a metric name, are working as expected.
func TestPrometheusUsesOptions(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, sampleTextFormat)
	}))
	defer ts.Close()

	p := &Prometheus{
		Urls:         []string{ts.URL},
		DotSeparator: true,
		MetricPrefix: "mytest",
	}

	var acc testutil.Accumulator
	// acc.SetDebug(true)

	err := p.Gather(&acc)
	require.NoError(t, err)

	assert.True(t, acc.HasFloatField("mytest.test.classes.total", "gauge"))
}
