package replication_test

import (
	"testing"
//	"time"

	"github.com/swapbyt3s/zenit/config"
//	"github.com/swapbyt3s/zenit/plugins/alerts/mysql/replication"
//	"github.com/swapbyt3s/zenit/plugins/inputs/percona/delay"
//	"github.com/swapbyt3s/zenit/plugins/lists/checks"
//	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

func TestMain(m *testing.M) {
	// Configure:
	config.File.MySQL.Alerts.Replication.Enable   = true
	config.File.MySQL.Alerts.Replication.Duration = 4
	config.File.MySQL.Inputs.Slave = true
	config.File.Process.Inputs.PerconaToolKitSlaveDelay = true

	// Run Tests:
	m.Run()
}

func TestConnection(t *testing.T) {
//	var histogram = []struct{
//		IO uint64
//		SQL uint64
//		Error uint64
//		Delay uint64
//		Status uint8
//		Notify bool
//	}{
//		{ 1, 0, 0, 0, checks.Normal   , false }, // 0s
//		{ 0, 0, 0, 0, checks.Normal   , false }, // 1s
//		{ 0, 0, 0, 0, checks.Normal   , false }, // 2s
//		{ 0, 0, 0, 0, checks.Normal   , false }, // 3s
//		{ 0, 0, 0, 0, checks.Normal   , false }, // 4s
//		{ 0, 0, 0, 0, checks.Normal   , false }, // 5s
//		{ 0, 0, 0, 0, checks.Notified , true  }, // 6s
//		{ 0, 0, 0, 0, checks.Recovered, true  }, // 7s
//	}
//
//	for second, variable := range histogram {
//
//		// Wait:
//		time.Sleep(1 * time.Second)
//	}
}
