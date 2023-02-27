package process

import (
	"github.com/debeando/zenit/common"
	"github.com/debeando/zenit/common/log"
	"github.com/debeando/zenit/config"
	"github.com/debeando/zenit/plugins/inputs"
	"github.com/debeando/zenit/plugins/lists/metrics"
)

type InputsPerconaKill struct{}

func (l *InputsPerconaKill) Collect() {
	defer func() {
		if err := recover(); err != nil {
			log.Error("InputsPerconaKill", map[string]interface{}{"error": err})
		}
	}()

	if !config.File.Inputs.Process.PerconaToolKitKill {
		return
	}

	var a = metrics.Load()
	var pid = common.PGrep("pt-kill")
	var value = 0

	if pid > 0 {
		value = 1
	}

	log.Debug("InputsPerconaKill", map[string]interface{}{"pt_kill": value})

	a.Add(metrics.Metric{
		Key: "process_pt_kill",
		Tags: []metrics.Tag{
			{Name: "hostname", Value: config.File.General.Hostname},
		},
		Values: []metrics.Value{
			{Key: "pt_kill", Value: value},
		},
	})
}

func init() {
	inputs.Add("InputsPerconaKill", func() inputs.Input { return &InputsPerconaKill{} })
}
