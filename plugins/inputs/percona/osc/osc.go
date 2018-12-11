// TODO: Read from zenit.yaml the list of process to check.
package process

import (
	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/loader"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

type InputsPerconaOSC struct {}

func (l *InputsPerconaOSC) Collect() {
	if ! config.File.Process.Inputs.PerconaToolKitOnlineSchemaChange {
		return
	}

	metrics.Load().Add(metrics.Metric{
		Key: "zenit_os",
		Tags: []metrics.Tag{{"system", "linux"},
			{"process", "pt_online_schema_change"}},
		Values: common.PGrep("pt-online-schema-change") ^ 1,
	})
}

func init() {
	loader.Add("InputsPerconaOSC", func() loader.Plugin { return &InputsPerconaOSC{} })
}
