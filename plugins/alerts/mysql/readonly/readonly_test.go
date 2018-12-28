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
		ReadOnly uint64
		Status uint8
	}{
		{ 1, 0 },
		{ 0, 0 },
		{ 0, 2 },
		{ 1, 2 },
		{ 1, 3 },
	}

	for _, check := range checks {
		// Add test value on metrics:
		metrics.Load().Reset()
		metrics.Load().Add(metrics.Metric{
			Key: "zenit_mysql_variables",
			Tags: []metrics.Tag{
				{"name", "read_only"},
			},
			Values: check.ReadOnly,
		})

		// Register alert:
		var c readonly.MySQLReadOnly
		c.Collect()

		// Evaluate alert status
		alert := alerts.Load().Exist("readonly")
		alert.Evaluate()

		if alert.Status != check.Status {
			t.Errorf("Expected: '%d', got: '%d'.", check.Status, alert.Status)
		}

		// Wait:
		time.Sleep(2 * time.Second)
	}
}
