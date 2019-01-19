package connections_test

import (
	"testing"
	"time"

	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/alerts/mysql/connections"
	"github.com/swapbyt3s/zenit/plugins/lists/checks"
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
// ASCII Graph for connections:
//
// MAX 100 ┤---------------------------------------------------------------------------------
//      90 ┤
//      80 ┤
//      70 ┤
//      60 ┤        ╭───╮   ╭─────────────w─╮
//      50 ┤        │   │   │               │
//      40 ┤        │   │   │               │
//      30 ┤    ╭───╯   │   │               │
//      20 ┤    │       │   │               │
//      10 ┤────╯       ╰───╯               ╰─c─
//       0 ---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|--
//        0  01  02  03  04  05  06  07  08  09  10  11  12  13  14  15  16  17  18  19  20
//
// Legend:
// - X = Time (s)
// - Y = Percentage (%)
// - w = Notify Warning
// - c = Notify Critical
// - r = Notify Recovered
//
	var histogram = []struct{
		MaxConnections uint64
		ThreadsConnected uint64
		Status uint8
		Notify bool
	}{
		{ 100, 10, checks.Normal   , false }, // 0s
		{ 100, 30, checks.Normal   , false }, // 1s
		{ 100, 60, checks.Normal   , false }, // 2s
		{ 100, 10, checks.Normal   , false }, // 3s
		{ 100, 60, checks.Normal   , false }, // 4s
		{ 100, 60, checks.Normal   , false }, // 5s
		{ 100, 60, checks.Normal   , false }, // 6s
		{ 100, 60, checks.Notified , true  }, // 7s
		{ 100, 10, checks.Recovered, true  }, // 8s
	}

	for second, variable := range histogram {
		// Add test value on metrics:
		metrics.Load().Reset()
		metrics.Load().Add(metrics.Metric{
			Key: "zenit_mysql_variables",
			Tags: []metrics.Tag{
				{"name", "max_connections"},
			},
			Values: variable.MaxConnections,
		})
		metrics.Load().Add(metrics.Metric{
			Key: "zenit_mysql_status",
			Tags: []metrics.Tag{
				{"name", "Threads_connected"},
			},
			Values: variable.ThreadsConnected,
		})

		// Register alert:
		var c connections.MySQLConnections
		c.Collect()

		check := checks.Load().Exist("connections")
		notify := check.Notify()

		if ! (check.Status == variable.Status && variable.Notify == notify) {
			t.Errorf("Second: %d, ThreadsConnected: %d, Evaluated: %t, Expected: %d, Got: %d.",
				second,
				variable.ThreadsConnected,
				notify,
				variable.Status,
				check.Status,
			)
		}

		// Wait:
		time.Sleep(1 * time.Second)
	}
}
