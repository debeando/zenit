package lagging_test

import (
	"testing"
	"time"

	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/alerts/mysql/lagging"
	"github.com/swapbyt3s/zenit/plugins/lists/checks"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

func TestMain(m *testing.M) {
	// Configure:
	config.File.MySQL.Alerts.Replication.Critical = 60
	config.File.MySQL.Alerts.Replication.Duration = 4
	config.File.MySQL.Alerts.Replication.Enable   = true
	config.File.MySQL.Alerts.Replication.Warning  = 10
	config.File.MySQL.Inputs.Slave = true

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
//      10 ┤    │       ╰───╯               ╰─c─
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
		Value uint64
		Status uint8
		Notify bool
	}{
		{  0, checks.Normal   , false }, // 0s
		{  0, checks.Normal   , false }, // 1s
		{ 10, checks.Normal   , false }, // 2s
		{ 10, checks.Normal   , false }, // 3s
		{ 10, checks.Normal   , false }, // 4s
		{  0, checks.Normal   , false }, // 5s
		{  0, checks.Normal   , false }, // 6s
		{  0, checks.Normal   , false }, // 7s
		{ 60, checks.Normal   , false }, // 8s
		{ 70, checks.Normal   , false }, // 9s
		{ 80, checks.Normal   , false }, // 10s
		{ 80, checks.Notified , true  }, // 11s
		{ 80, checks.Notified , false }, // 12s
		{  0, checks.Recovered, true  }, // 13s
		{  0, checks.Normal   , false }, // 14s
	}

	for second, variable := range histogram {
		// Add test value on metrics:
		metrics.Load().Reset()
		metrics.Load().Add(metrics.Metric{
			Key: "zenit_mysql_slave",
			Tags: []metrics.Tag{
				{"name", "Seconds_Behind_Master"},
			},
			Values: variable.Value,
		})

		// Register alert:
		var c lagging.MySQLLagging
		c.Collect()

		// Evaluate alert status
		check := checks.Load().Exist("lagging")
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
