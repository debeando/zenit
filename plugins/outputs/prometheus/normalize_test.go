package prometheus_test

import (
	"testing"

	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
	"github.com/swapbyt3s/zenit/plugins/outputs/prometheus"
)

var a = metrics.Load()

func TestNormalize(t *testing.T) {
	a.Add(metrics.Metric{
		Key: "test_metric",
		Tags: []metrics.Tag{
			{"foo", "bar"},
		},
		Values: []metrics.Value{
			{"a", uint(2)},
			{"b", uint(2)},
		},
	})

	expected := "test_metric{foo=\"bar\",type=\"a\"} 2\ntest_metric{foo=\"bar\",type=\"b\"} 2\n"
	result := prometheus.Normalize(a)

	if result != expected {
		t.Errorf("Expected: '%s', got: '%s'", expected, result)
	}
}
