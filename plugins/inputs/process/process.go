// TODO: Read from zenit.yaml the list of process to check.
package process

import (
	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/plugins/accumulator"
)

func PerconaToolKitKill() {
	accumulator.Load().Add(accumulator.Metric{
		Key: "os",
		Tags: []accumulator.Tag{{"system", "linux"},
			{"process", "pt_kill"}},
		Values: common.PGrep("pt-kill") ^ 1,
	})
}

func PerconaToolKitDeadlockLogger() {
	accumulator.Load().Add(accumulator.Metric{
		Key: "os",
		Tags: []accumulator.Tag{{"system", "linux"},
			{"process", "pt_deadlock_logger"}},
		Values: common.PGrep("pt-deadlock-logger") ^ 1,
	})
}

func PerconaToolKitSlaveDelay() {
	accumulator.Load().Add(accumulator.Metric{
		Key: "os",
		Tags: []accumulator.Tag{{"system", "linux"},
			{"process", "pt_slave_delay"}},
		Values: common.PGrep("pt-slave-delay") ^ 1,
	})
}
