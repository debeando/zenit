package connections_test

import (
	"testing"
	"time"

	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/alerts/mysql/connections"
	"github.com/swapbyt3s/zenit/plugins/lists/alerts"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

func TestMain(m *testing.M) {
	// Configure:
	config.File.MySQL.Alerts.Connections.Critical = 70
	config.File.MySQL.Alerts.Connections.Duration = 4
	config.File.MySQL.Alerts.Connections.Enable   = true
	config.File.MySQL.Alerts.Connections.Warning  = 50
	config.File.MySQL.Inputs.Status = true
	config.File.MySQL.Inputs.Variables = true

	// Run Tests:
	m.Run()
}

func TestConnection(t *testing.T) {
	var checks = []struct{
		MaxConnections uint64
		ThreadsConnected uint64
		Status uint8
	}{
		{ 100, 10, 0 },
		{ 100, 60, 0 },
		{ 100, 60, 1 },
		{ 100, 70, 1 },
		{ 100, 70, 2 },
		{ 100, 30, 2 },
		{ 100, 30, 3 },
		{ 100, 90, 3 },
		{ 100, 90, 2 },
		{ 100, 10, 2 },
		{ 100, 10, 3 },
	}

	for _, check := range checks {
		// Add test value on metrics:
		metrics.Load().Reset()
		metrics.Load().Add(metrics.Metric{
			Key: "zenit_mysql_variables",
			Tags: []metrics.Tag{
				{"name", "max_connections"},
			},
			Values: check.MaxConnections,
		})
		metrics.Load().Add(metrics.Metric{
			Key: "zenit_mysql_status",
			Tags: []metrics.Tag{
				{"name", "Threads_connected"},
			},
			Values: check.ThreadsConnected,
		})

		// Register alert:
		var c connections.MySQLConnections
		c.Collect()

		// Evaluate alert status
		alert := alerts.Load().Exist("connections")
		alert.Evaluate()

		if alert.Status != check.Status {
			t.Errorf("Expected: '%d', got: '%d'.", check.Status, alert.Status)
		}

		// Wait:
		time.Sleep(2 * time.Second)
	}
}
