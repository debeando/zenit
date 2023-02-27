package process

import (
	"github.com/debeando/zenit/common"
	"github.com/debeando/zenit/common/log"
	"github.com/debeando/zenit/config"
	"github.com/debeando/zenit/plugins/inputs"
	"github.com/debeando/zenit/plugins/lists/metrics"
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
			{Name: "hostname", Value: config.File.General.Hostname},
		},
		Values: []metrics.Value{
			{Key: "pt_online_schema_change", Value: value},
		},
	})
}

func init() {
	inputs.Add("InputsPerconaOSC", func() inputs.Input { return &InputsPerconaOSC{} })
}
