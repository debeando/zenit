package process

import (
	"zenit/agent/plugins/inputs"
	"zenit/agent/plugins/lists/metrics"
	"zenit/config"

	"github.com/debeando/go-common/exec"
	"github.com/debeando/go-common/log"
)

type Plugin struct{}

func (p *Plugin) Collect(name string, cnf *config.Config, mtc *metrics.Items) {
	defer func() {
		if err := recover(); err != nil {
			log.ErrorWithFields(name, log.Fields{"message": err})
		}
	}()

	if !cnf.Inputs.Process.PerconaToolKitOnlineSchemaChange {
		return
	}

	var pid = exec.PGrep("pt-online-schema-change")
	var value = 0

	if pid > 0 {
		value = 1
	}

	log.DebugWithFields(name, log.Fields{"pt_online_schema_change": value})

	mtc.Add(metrics.Metric{
		Key: "process_pt_online_schema_change",
		Tags: []metrics.Tag{
			{Name: "hostname", Value: cnf.General.Hostname},
		},
		Values: []metrics.Value{
			{Key: "pt_online_schema_change", Value: value},
		},
	})
}

func init() {
	inputs.Add("InputsPerconaOSC", func() inputs.Input { return &Plugin{} })
}
