// TODO: Read from zenit.yaml the list of process to check.
package process

import (
	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/loader"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

type InputsPerconaDeadlock struct {}

func (l *InputsPerconaDeadlock) Collect() {
	if ! config.File.Process.Inputs.PerconaToolKitDeadlockLogger {
		return
	}

	metrics.Load().Add(metrics.Metric{
		Key: "zenit_os",
		Tags: []metrics.Tag{{"system", "linux"},
			{"process", "pt_deadlock_logger"}},
		Values: common.PGrep("pt-deadlock-logger") ^ 1,
	})
}

func init() {
	loader.Add("InputsPerconaDeadlock", func() loader.Plugin { return &InputsPerconaDeadlock{} })
}
