package prometheus_test

import (
  "strings"
  "testing"

  "github.com/swapbyt3s/zenit/plugins/lists/accumulator"
  "github.com/swapbyt3s/zenit/plugins/outputs/prometheus"
)

var a = accumulator.Load()

func TestNormalize(t *testing.T) {
  a.Add(accumulator.Metric{
    Key:  "test_metric",
    Tags: []accumulator.Tag{
      {"foo", "bar"},
    },
    Values: []accumulator.Value{
      {"a", uint(2)},
      {"b", uint(2)},
    },
  })

  expected := "test_metric{foo=\"bar\",type=\"a\"} 2\ntest_metric{foo=\"bar\",type=\"b\"} 2"
  output := prometheus.Normalize(a)
  result := strings.Join(output, "\n")

  if result != expected {
    t.Error("Expected: " + expected)
  }
}
