package process

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
	"github.com/swapbyt3s/zenit/plugins/inputs"
)

type InputsPerconaDelay struct {}

func (l *InputsPerconaDelay) Collect() {
	if ! config.File.Process.Inputs.PerconaToolKitSlaveDelay {
		return
	}

	var pid = common.PGrep("pt-slave-delay")
	var value = 0

	if pid > 0 {
		value = 1
	}

	metrics.Load().Add(metrics.Metric{
		Key: "zenit_process",
		Tags: []metrics.Tag{
			{"system", "linux"},
			{"name", "pt_slave_delay"},
		},
		Values: value,
	})

	log.Debug(fmt.Sprintf("Plugin - InputsPerconaDelay - %d", value))
}

func init() {
	inputs.Add("InputsPerconaDelay", func() inputs.Input { return &InputsPerconaDelay{} })
}
