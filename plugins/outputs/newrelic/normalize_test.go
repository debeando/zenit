package newrelic_test

import (
	"testing"

	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
	"github.com/swapbyt3s/zenit/plugins/outputs/newrelic"
)

var a = metrics.Load()

func TestMain(t *testing.T) {
	config.File.General.Hostname = "localhost.test"
}

func TestNormalizeSingleValue(t *testing.T) {
	a.Add(metrics.Metric{
		Key: "test_metric",
		Tags: []metrics.Tag{
			{"name", "bar"},
		},
		Values: uint(1),
	})

	result := newrelic.Normalize(a)

	if _, ok := result["test_metric"]["eventType"]; ! ok {
		t.Errorf("Expected key: eventType, got: '%s'", result)
	}

	if _, ok := result["test_metric"]["hostname"]; ! ok {
		t.Errorf("Expected key: hostname, got: '%s'", result)
	}

	if _, ok := result["test_metric"]["bar"]; ! ok {
		t.Errorf("Expected key: bar, got: '%s'", result)
	}

	if value := result["test_metric"]["bar"]; value != uint(1) {
		t.Errorf("Expected key: bar, got: '%s'", result)
	}
}

func TestNormalizeMultipleValues(t *testing.T) {
	a.Add(metrics.Metric{
		Key: "test_metric",
		Tags: []metrics.Tag{
			{"foo", "bar"},
		},
		Values: []metrics.Value{
			{"a", uint(1)},
			{"b", uint(2)},
			{"c", uint(3)},
		},
	})

	result := newrelic.Normalize(a)

	if _, ok := result["test_metric"]["foo"]; ! ok {
		t.Errorf("Expected key: foo, got: '%s'", result)
	}

	if value := result["test_metric"]["foo"]; value != "bar" {
		t.Errorf("Expected key: foo, got: '%s'", result)
	}

	if _, ok := result["test_metric"]["a"]; ! ok {
		t.Errorf("Expected key: a, got: '%s'", result)
	}

	if value := result["test_metric"]["a"]; value != uint(1) {
		t.Errorf("Expected key: a, got: '%s'", result)
	}
}
