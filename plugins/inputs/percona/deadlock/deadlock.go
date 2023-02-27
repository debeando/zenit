package process

import (
	"github.com/debeando/zenit/common"
	"github.com/debeando/zenit/common/log"
	"github.com/debeando/zenit/config"
	"github.com/debeando/zenit/plugins/inputs"
	"github.com/debeando/zenit/plugins/lists/metrics"
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
			{Name: "hostname", Value: config.File.General.Hostname},
		},
		Values: []metrics.Value{
			{Key: "pt_deadlock_logger", Value: value},
		},
	})
}

func init() {
	inputs.Add("InputsPerconaDeadlock", func() inputs.Input { return &InputsPerconaDeadlock{} })
}
