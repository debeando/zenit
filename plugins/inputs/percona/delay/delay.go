package process

import (
	"github.com/debeando/zenit/common"
	"github.com/debeando/zenit/common/log"
	"github.com/debeando/zenit/config"
	"github.com/debeando/zenit/plugins/inputs"
	"github.com/debeando/zenit/plugins/lists/metrics"
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
			{Name: "hostname", Value: config.File.General.Hostname},
		},
		Values: []metrics.Value{
			{Key: "pt_slave_delay", Value: value},
		},
	})
}

func init() {
	inputs.Add("InputsPerconaToolkitSlaveDelay", func() inputs.Input { return &InputsPerconaToolkitSlaveDelay{} })
}
