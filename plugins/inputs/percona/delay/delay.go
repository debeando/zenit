package process

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
	"github.com/swapbyt3s/zenit/plugins/inputs"
)

type InputsPerconaToolkitSlaveDelay struct {}

func (l *InputsPerconaToolkitSlaveDelay) Collect() {
	defer func () {
		if err := recover(); err != nil {
			log.Debug(fmt.Sprintf("Plugin - InputsPerconaToolkitSlaveDelay - Panic (code %d) has been recover from somewhere.\n", err))
		}
	}()

	if ! config.File.Process.Inputs.PerconaToolKitSlaveDelay {
		return
	}

	var pid = common.PGrep("pt-slave-delay")
	var value uint64 = 0

	if pid > 0 {
		value = 1
	}

	metrics.Load().Add(metrics.Metric{
		Key: "zenit_process",
		Tags: []metrics.Tag{
			{"name", "pt_slave_delay"},
		},
		Values: value,
	})

	log.Debug(fmt.Sprintf("Plugin - InputsPerconaToolkitSlaveDelay - %d", value))
}

func init() {
	inputs.Add("InputsPerconaToolkitSlaveDelay", func() inputs.Input { return &InputsPerconaToolkitSlaveDelay{} })
}
