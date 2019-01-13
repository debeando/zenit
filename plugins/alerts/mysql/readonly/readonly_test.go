package readonly_test

import (
	"testing"
	"time"

	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/alerts/mysql/readonly"
	"github.com/swapbyt3s/zenit/plugins/lists/checks"
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
	var histogram = []struct{
		Value uint64
		Status uint8
		Notify bool
	}{
		{ 1, checks.Normal   , false }, // 0s
		{ 0, checks.Normal   , false }, // 1s
		{ 1, checks.Normal   , false }, // 2s
		{ 0, checks.Normal   , false }, // 3s
		{ 0, checks.Normal   , false }, // 4s
		{ 0, checks.Normal   , false }, // 5s
		{ 0, checks.Notified , true  }, // 6s
		{ 1, checks.Recovered, true  }, // 7s
	}

	for second, variable := range histogram {
		// Add test value on metrics:
		metrics.Load().Reset()
		metrics.Load().Add(metrics.Metric{
			Key: "zenit_mysql_variables",
			Tags: []metrics.Tag{
				{"name", "read_only"},
			},
			Values: variable.Value,
		})

		// Register alert:
		var c readonly.MySQLReadOnly
		c.Collect()

		// Evaluate alert status
		check := checks.Load().Exist("readonly")
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
