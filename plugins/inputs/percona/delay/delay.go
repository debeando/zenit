// TODO: Read from zenit.yaml the list of process to check.
package process

import (
	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/loader"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

type InputsPerconaDelay struct {}

func (l *InputsPerconaDelay) Collect() {
	if ! config.File.Process.Inputs.PerconaToolKitSlaveDelay {
		return
	}

	metrics.Load().Add(metrics.Metric{
		Key: "zenit_os",
		Tags: []metrics.Tag{{"system", "linux"},
			{"process", "pt_slave_delay"}},
		Values: common.PGrep("pt-slave-delay") ^ 1,
	})
}

func init() {
	loader.Add("InputsPerconaDelay", func() loader.Plugin { return &InputsPerconaDelay{} })
}
