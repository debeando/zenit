// TODO: Read from zenit.yaml the list of process to check.
package process

import (
	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/loader"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

type InputsPerconaKill struct {}

func (l *InputsPerconaKill) Collect() {
	if ! config.File.Process.Inputs.PerconaToolKitKill {
		return
	}

	metrics.Load().Add(metrics.Metric{
		Key: "zenit_os",
		Tags: []metrics.Tag{{"system", "linux"},
			{"process", "pt_kill"}},
		Values: common.PGrep("pt-kill") ^ 1,
	})
}

func init() {
	loader.Add("InputsPerconaKill", func() loader.Plugin { return &InputsPerconaKill{} })
}
