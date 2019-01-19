package process

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
	"github.com/swapbyt3s/zenit/plugins/inputs"
)

type InputsPerconaXtraBackup struct {}

func (l *InputsPerconaXtraBackup) Collect() {
	if ! config.File.Process.Inputs.PerconaXtraBackup {
		return
	}

	var pid = common.PGrep("xtrabackup")
	var value = 0

	if pid > 0 {
		value = 1
	}

	metrics.Load().Add(metrics.Metric{
		Key: "zenit_process",
		Tags: []metrics.Tag{
			{"name", "xtrabackup"},
		},
		Values: value,
	})

	log.Debug(fmt.Sprintf("Plugin - InputsPerconaOSC - %d", value))
}

func init() {
	inputs.Add("InputsPerconaXtraBackup", func() inputs.Input { return &InputsPerconaXtraBackup{} })
}
