package process

import (
	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/inputs"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

type InputsPerconaDeadlock struct{}

func (l *InputsPerconaDeadlock) Collect() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("InputsPerconaDeadlock", map[string]interface{}{"error": err})
		}
	}()

	if !config.File.Inputs.Process.PerconaToolKitDeadlockLogger {
		return
	}

	var a = metrics.Load()
	var pid = common.PGrep("pt-deadlock-logger")
	var value = 0

	if pid > 0 {
		value = 1
	}

	log.Debug("InputsPerconaDeadlock", map[string]interface{}{"pt_deadlock_logger": value})

	a.Add(metrics.Metric{
		Key: "process_pt_deadlock_logger",
		Tags: []metrics.Tag{
			{"hostname", config.File.General.Hostname},
		},
		Values: []metrics.Value{
			{"pt_deadlock_logger", value},
		},
	})
}

func init() {
	inputs.Add("InputsPerconaDeadlock", func() inputs.Input { return &InputsPerconaDeadlock{} })
}
