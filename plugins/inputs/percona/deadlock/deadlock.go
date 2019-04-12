package process

import (
	"fmt"

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
			log.Debug(fmt.Sprintf("Plugin - InputsPerconaDeadlock - Panic (code %d) has been recover from somewhere.\n", err))
		}
	}()

	if !config.File.Process.Inputs.PerconaToolKitDeadlockLogger {
		return
	}

	var pid = common.PGrep("pt-deadlock-logger")
	var value = 0

	if pid > 0 {
		value = 1
	}

	metrics.Load().Add(metrics.Metric{
		Key: "zenit_process",
		Tags: []metrics.Tag{
			{"name", "pt_deadlock_logger"},
		},
		Values: value,
	})

	log.Debug(fmt.Sprintf("Plugin - InputsPerconaDeadlock - %d", value))
}

func init() {
	inputs.Add("InputsPerconaDeadlock", func() inputs.Input { return &InputsPerconaDeadlock{} })
}
