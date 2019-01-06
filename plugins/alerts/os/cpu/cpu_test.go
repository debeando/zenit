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
		Value uint64
		Status uint8
		Notify bool
	}{
		{ 10, alerts.Normal   , false }, // 1s
		{ 10, alerts.Normal   , false }, // 2s
		{ 50, alerts.Normal   , false }, // 3s
		{ 10, alerts.Normal   , false }, // 4s
		{ 10, alerts.Normal   , false }, // 5s
		{ 50, alerts.Normal   , false }, // 6s
		{ 55, alerts.Normal   , false }, // 7s
		{ 50, alerts.Normal   , false }, // 8s
		{ 90, alerts.Notified , true  }, // 9s
		{ 95, alerts.Normal   , false }, // 10s
		{ 10, alerts.Recovered, true  }, // 11s
		{ 10, alerts.Normal   , false }, // 12s
	}

	for second, check := range checks {
		// Add test value on metrics:
		metrics.Load().Reset()
		metrics.Load().Add(metrics.Metric{
			Key: "zenit_os",
			Tags: []metrics.Tag{
				{"name", "cpu"},
			},
			Values: check.Value,
		})

		// Register alert:
		var c cpu.OSCPU
		c.Collect()

		// Evaluate alert status
		alert := alerts.Load().Exist("cpu")
		notify := alert.Notify()

		if ! (alert.Status == check.Status && check.Notify == notify) {
			t.Errorf("Second: %d, Value: %d, Evaluated: %t, Expected: '%d', Got: '%d'.",
				second,
				check.Value,
				notify,
				check.Status,
				alert.Status,
			)
		}

		// Wait:
		time.Sleep(1 * time.Second)
	}
}
