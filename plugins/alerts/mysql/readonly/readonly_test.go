package readonly_test

import (
	"testing"
	"time"

	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/alerts/mysql/readonly"
	"github.com/swapbyt3s/zenit/plugins/lists/alerts"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

func TestMain(m *testing.M) {
	// Configure:
	config.File.MySQL.Alerts.ReadOnly.Enable   = true
	config.File.MySQL.Alerts.ReadOnly.Duration = 4
	config.File.MySQL.Inputs.Variables = true

	// Run Tests:
	m.Run()
}

func TestConnection(t *testing.T) {
	var checks = []struct{
		Value uint64
		Status uint8
		Notify bool
	}{
		{ 1, alerts.Normal   , false }, // 0s
		{ 0, alerts.Normal   , false }, // 1s
		{ 1, alerts.Normal   , false }, // 2s
		{ 0, alerts.Normal   , false }, // 3s
		{ 0, alerts.Normal   , false }, // 4s
		{ 0, alerts.Normal   , false }, // 5s
		{ 0, alerts.Notified , true  }, // 6s
		{ 1, alerts.Recovered, true  }, // 7s
	}

	for second, check := range checks {
		// Add test value on metrics:
		metrics.Load().Reset()
		metrics.Load().Add(metrics.Metric{
			Key: "zenit_mysql_variables",
			Tags: []metrics.Tag{
				{"name", "read_only"},
			},
			Values: check.Value,
		})

		// Register alert:
		var c readonly.MySQLReadOnly
		c.Collect()

		// Evaluate alert status
		alert := alerts.Load().Exist("readonly")
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
