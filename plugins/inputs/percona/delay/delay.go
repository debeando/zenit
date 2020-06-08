package process

import (
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
			log.Error("InputsPerconaToolkitSlaveDelay", map[string]interface{}{"error": err})
		}
	}()

	if !config.File.Inputs.Process.PerconaToolKitSlaveDelay {
		return
	}

	var a = metrics.Load()
	var pid = common.PGrep("pt-slave-delay")
	var value int64 = 0

	if pid > 0 {
		value = 1
	}

	log.Debug("InputsPerconaToolkitSlaveDelay", map[string]interface{}{"pt_slave_delay": value})

	a.Add(metrics.Metric{
		Key: "process_pt_slave_delay",
		Tags: []metrics.Tag{
			{"hostname", config.File.General.Hostname},
		},
		Values: []metrics.Value{
			{ "pt_slave_delay", value},
		},
	})
}

func init() {
	inputs.Add("InputsPerconaToolkitSlaveDelay", func() inputs.Input { return &InputsPerconaToolkitSlaveDelay{} })
}
