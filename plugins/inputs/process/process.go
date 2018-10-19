// TODO: Read from zenit.yaml the list of process to check.
package process

import (
	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

func PerconaToolKitKill() {
	metrics.Load().Add(metrics.Metric{
		Key: "os",
		Tags: []metrics.Tag{{"system", "linux"},
			{"process", "pt_kill"}},
		Values: common.PGrep("pt-kill") ^ 1,
	})
}

func PerconaToolKitDeadlockLogger() {
	metrics.Load().Add(metrics.Metric{
		Key: "os",
		Tags: []metrics.Tag{{"system", "linux"},
			{"process", "pt_deadlock_logger"}},
		Values: common.PGrep("pt-deadlock-logger") ^ 1,
	})
}

func PerconaToolKitSlaveDelay() {
	metrics.Load().Add(metrics.Metric{
		Key: "os",
		Tags: []metrics.Tag{{"system", "linux"},
			{"process", "pt_slave_delay"}},
		Values: common.PGrep("pt-slave-delay") ^ 1,
	})
}

func PerconaToolKitOnlineSchemaChange() {
	metrics.Load().Add(metrics.Metric{
		Key: "os",
		Tags: []metrics.Tag{{"system", "linux"},
			{"process", "pt_online_schema_change"}},
		Values: common.PGrep("pt-online-schema-change") ^ 1,
	})
}

func PerconaXtraBackup() {
	metrics.Load().Add(metrics.Metric{
		Key: "os",
		Tags: []metrics.Tag{{"system", "linux"},
			{"process", "xtrabackup"}},
		Values: common.PGrep("xtrabackup") ^ 1,
	})
}
