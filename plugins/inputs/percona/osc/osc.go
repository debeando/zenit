package process

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/loader"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

type InputsPerconaOSC struct {}

func (l *InputsPerconaOSC) Collect() {
	if ! config.File.Process.Inputs.PerconaToolKitOnlineSchemaChange {
		return
	}

	var pid = common.PGrep("pt-online-schema-change")
	var value = 0

	if pid > 0 {
		value = 1
	}

	metrics.Load().Add(metrics.Metric{
		Key: "zenit_process",
		Tags: []metrics.Tag{
			{"system", "linux"},
			{"name", "pt_online_schema_change"},
		},
		Values: value,
	})

	log.Debug(fmt.Sprintf("Plugin - InputsPerconaOSC - %d", value))
}

func init() {
	loader.Add("InputsPerconaOSC", func() loader.Plugin { return &InputsPerconaOSC{} })
}
