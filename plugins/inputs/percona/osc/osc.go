package process

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/inputs"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

type InputsPerconaOSC struct{}

func (l *InputsPerconaOSC) Collect() {
	defer func() {
		if err := recover(); err != nil {
			log.Debug(fmt.Sprintf("Plugin - InputsPerconaOSC - Panic (code %d) has been recover from somewhere.\n", err))
		}
	}()

	if !config.File.Inputs.Process.PerconaToolKitOnlineSchemaChange {
		return
	}

	var a = metrics.Load()
	var pid = common.PGrep("pt-online-schema-change")
	var value = 0

	if pid > 0 {
		value = 1
	}

	a.Add(metrics.Metric{
		Key: "process",
		Tags: []metrics.Tag{
			{"hostname", config.File.General.Hostname},
		},
		Values: []metrics.Value{
			{ "pt_online_schema_change", value},
		},
	})

	log.Debug(fmt.Sprintf("Plugin - InputsPerconaOSC - %d", value))
}

func init() {
	inputs.Add("InputsPerconaOSC", func() inputs.Input { return &InputsPerconaOSC{} })
}
