package cpu_test

import (
	"testing"
	"time"

	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/alerts/os/cpu"
	"github.com/swapbyt3s/zenit/plugins/lists/checks"
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
	var histogram = []struct{
		Value uint64
		Status uint8
		Notify bool
	}{
		{ 10, checks.Normal   , false }, // 1s
		{ 10, checks.Normal   , false }, // 2s
		{ 50, checks.Normal   , false }, // 3s
		{ 10, checks.Normal   , false }, // 4s
		{ 10, checks.Normal   , false }, // 5s
		{ 50, checks.Normal   , false }, // 6s
		{ 55, checks.Normal   , false }, // 7s
		{ 50, checks.Normal   , false }, // 8s
		{ 90, checks.Notified , true  }, // 9s
		{ 95, checks.Normal   , false }, // 10s
		{ 10, checks.Recovered, true  }, // 11s
		{ 10, checks.Normal   , false }, // 12s
	}

	for second, variable := range histogram {
		// Add test value on metrics:
		metrics.Load().Reset()
		metrics.Load().Add(metrics.Metric{
			Key: "zenit_os",
			Tags: []metrics.Tag{
				{"name", "cpu"},
			},
			Values: variable.Value,
		})

		// Register alert:
		var c cpu.OSCPU
		c.Collect()

		// Evaluate alert status
		check := checks.Load().Exist("cpu")
		notify := check.Notify()

		if ! (check.Status == variable.Status && variable.Notify == notify) {
			t.Errorf("Second: %d, Value: %d, Evaluated: %t, Expected: '%d', Got: '%d'.",
				second,
				variable.Value,
				notify,
				variable.Status,
				check.Status,
			)
		}

		// Wait:
		time.Sleep(1 * time.Second)
	}
}
