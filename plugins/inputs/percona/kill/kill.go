package process

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
	"github.com/swapbyt3s/zenit/plugins/inputs"
)

type InputsPerconaKill struct {}

func (l *InputsPerconaKill) Collect() {
	if ! config.File.Process.Inputs.PerconaToolKitKill {
		return
	}

	var pid = common.PGrep("pt-kill")
	var value = 0

	if pid > 0 {
		value = 1
	}

	metrics.Load().Add(metrics.Metric{
		Key: "zenit_process",
		Tags: []metrics.Tag{
			{"system", "linux"},
			{"name", "pt_kill"},
		},
		Values: value,
	})

	log.Debug(fmt.Sprintf("Plugin - InputsPerconaKill - %d", value))
}

func init() {
	inputs.Add("InputsPerconaKill", func() inputs.Input { return &InputsPerconaKill{} })
}
