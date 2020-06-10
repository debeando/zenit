package process

import (
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
			log.Error("InputsPerconaOSC", map[string]interface{}{"error": err})
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

	log.Debug("InputsPerconaOSC", map[string]interface{}{"pt_online_schema_change": value})

	a.Add(metrics.Metric{
		Key: "process_pt_online_schema_change",
		Tags: []metrics.Tag{
			{"hostname", config.File.General.Hostname},
		},
		Values: []metrics.Value{
			{ "pt_online_schema_change", value},
		},
	})
}

func init() {
	inputs.Add("InputsPerconaOSC", func() inputs.Input { return &InputsPerconaOSC{} })
}
