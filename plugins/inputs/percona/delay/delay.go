package process

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/inputs"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

type InputsPerconaToolkitSlaveDelay struct{}

func (l *InputsPerconaToolkitSlaveDelay) Collect() {
	defer func() {
		if err := recover(); err != nil {
			log.Debug(fmt.Sprintf("Plugin - InputsPerconaToolkitSlaveDelay - Panic (code %d) has been recover from somewhere.\n", err))
		}
	}()

	if !config.File.Inputs.Process.PerconaToolKitSlaveDelay {
		return
	}

	log.Info("Plugin - InputsPerconaToolkitSlaveDelay")

	var a = metrics.Load()
	var pid = common.PGrep("pt-slave-delay")
	var value int64 = 0

	if pid > 0 {
		value = 1
	}

	a.Add(metrics.Metric{
		Key: "process_pt_slave_delay",
		Tags: []metrics.Tag{
			{"hostname", config.File.General.Hostname},
		},
		Values: []metrics.Value{
			{ "pt_slave_delay", value},
		},
	})

	log.Debug(fmt.Sprintf("Plugin - InputsPerconaToolkitSlaveDelay - %d", value))
}

func init() {
	inputs.Add("InputsPerconaToolkitSlaveDelay", func() inputs.Input { return &InputsPerconaToolkitSlaveDelay{} })
}
