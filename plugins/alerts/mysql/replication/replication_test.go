package replication_test

import (
	"testing"
	"time"

	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/alerts/mysql/replication"
	"github.com/swapbyt3s/zenit/plugins/lists/checks"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

func TestMain(m *testing.M) {
	// Configure:
	config.File.MySQL.Alerts.Replication.Enable   = true
	config.File.MySQL.Alerts.Replication.Duration = 3
	config.File.MySQL.Inputs.Slave = true
	config.File.Process.Inputs.PerconaToolKitSlaveDelay = true

	// Run Tests:
	m.Run()
}

func TestReplication(t *testing.T) {
// ASCII Graph for connections:
//
//    |-------------------------------------------------------------------------------
//    |
//    |  (IO)                (SQL)               (SQL & PTSD)        (Error & SQL)
//  1 ┤  ╭───────c───r       ╭───────c───r       ╭───────────╮       ╭───────c───r
//  0 ---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
//   0  01  02  03  04  05  06  07  08  09  10  11  12  13  14  15  16  17  18  19  20
//
// Legend:
// - X = Time (s)
// - Y = Boolean values (1 YES, 0 NO)
// - c = Notify Critical
// - r = Notify Recovered
//

	var histogram = []struct{
		IO uint64
		SQL uint64
		Error uint64
		PTSD uint64 // Percona Tool Slave Delay (pt-slave-delay)
		Status uint8
		Notify bool
	}{
		//I  S  E  P  Status            Notify
		{ 1, 1, 0, 0, checks.Normal   , false }, // 0s
		{ 0, 1, 0, 0, checks.Normal   , false }, // 1s
		{ 0, 1, 0, 0, checks.Normal   , false }, // 2s
		{ 0, 1, 0, 0, checks.Notified , true  }, // 3s
		{ 1, 1, 0, 0, checks.Recovered, true  }, // 4s
		{ 1, 1, 0, 0, checks.Normal   , false }, // 5s
		{ 1, 0, 0, 0, checks.Normal   , false }, // 6s
		{ 1, 0, 0, 0, checks.Normal   , false }, // 7s
		{ 1, 0, 0, 0, checks.Notified , true  }, // 8s
		{ 1, 1, 0, 0, checks.Recovered, true  }, // 9s
		{ 1, 1, 0, 0, checks.Normal   , false }, // 10s
		{ 1, 0, 0, 1, checks.Normal   , false }, // 11s
		{ 1, 0, 0, 1, checks.Normal   , false }, // 12s
		{ 1, 0, 0, 1, checks.Normal   , false }, // 13s
		{ 1, 1, 0, 0, checks.Normal   , false }, // 14s
		{ 1, 1, 0, 0, checks.Normal   , false }, // 15s
		{ 1, 0, 1, 0, checks.Normal   , false }, // 16s
		{ 1, 0, 1, 0, checks.Normal   , false }, // 17s
		{ 1, 0, 1, 0, checks.Notified , true  }, // 18s
		{ 1, 1, 0, 0, checks.Recovered, true  }, // 19s
		{ 1, 1, 0, 0, checks.Normal   , false }, // 20s
	}

	for second, variable := range histogram {
		// Add test value on metrics:
		metrics.Load().Reset()
		metrics.Load().Add(metrics.Metric{
			Key: "zenit_mysql_slave",
			Tags: []metrics.Tag{
				{"name", "Slave_IO_Running"},
			},
			Values: variable.IO,
		})
		metrics.Load().Add(metrics.Metric{
			Key: "zenit_mysql_slave",
			Tags: []metrics.Tag{
				{"name", "Slave_SQL_Running"},
			},
			Values: variable.SQL,
		})
		metrics.Load().Add(metrics.Metric{
			Key: "zenit_mysql_slave",
			Tags: []metrics.Tag{
				{"name", "Last_SQL_Errno"},
			},
			Values: variable.Error,
		})
		metrics.Load().Add(metrics.Metric{
			Key: "zenit_process",
			Tags: []metrics.Tag{
				{"name", "pt_slave_delay"},
			},
			Values: variable.PTSD,
		})

		// Register alert:
		var c replication.MySQLReplication
		c.Collect()

		// Evaluate alert status
		check := checks.Load().Exist("replication")
		notify := check.Notify()

		if ! (check.Status == variable.Status && variable.Notify == notify) {
			t.Errorf("Second: %d, IO: %d, SQL: %d, Error: %d, PTSD: %d, Evaluated: %t, Expected: %d, Got: %d.",
				second,
				variable.IO,
				variable.SQL,
				variable.Error,
				variable.PTSD,
				notify,
				variable.Status,
				check.Status,
			)
		}

		// Wait:
		time.Sleep(1 * time.Second)
	}
}
