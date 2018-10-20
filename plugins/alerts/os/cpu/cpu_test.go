package cpu_test

import (
	"testing"
	"time"

	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/alerts/os/cpu"
	"github.com/swapbyt3s/zenit/plugins/lists/alerts"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

func TestMain(m *testing.M) {
	// Configure:
	config.File.OS.Alerts.CPU.Enable   = true
	config.File.OS.Alerts.CPU.Duration = 4
	config.File.OS.Alerts.CPU.Warning  = 50
	config.File.OS.Alerts.CPU.Critical = 90

	// Run Tests:
	m.Run()
}

func TestCPU(t *testing.T) {
	var checks = []struct{
		Value  float64
		Status uint8
	}{
		{ 10, 0 }, // CPU Peak, False positive.
		{ 50, 0 }, //
		{ 10, 0 }, //
		{ 10, 0 }, //
		{ 50, 1 }, // Warning
		{ 55, 1 }, //
		{ 50, 1 }, //
		{ 90, 1 }, // Critical
		{ 95, 2 }, //
		{ 10, 2 }, //
		{ 10, 3 }, //
		{ 10, 3 }, // Resolved
	}

	for _, check := range checks {
		// Add test value on metrics:
		metrics.Load().Reset()
		metrics.Load().Add(metrics.Metric{
			Key: "os",
			Tags: []metrics.Tag{
				{"name", "cpu"},
			},
			Values: check.Value,
		})

		// Register alert:
		cpu.Register()

		// Evaluate alert status
		alert := alerts.Load().Exist("cpu")
		alert.Evaluate()

		if alert.Status != check.Status {
			t.Errorf("Expected: '%d', got: '%d'.", check.Status, alert.Status)
		}

		// Wait:
		time.Sleep(2 * time.Second)
	}
}
