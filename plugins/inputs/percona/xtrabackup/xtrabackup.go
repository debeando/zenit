package process

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/loader"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
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
		Key: "zenit_os",
		Tags: []metrics.Tag{
			{"system", "linux"},
			{"process", "xtrabackup"},
		},
		Values: value,
	})

	log.Debug(fmt.Sprintf("Plugin - InputsPerconaOSC - %d", value))
}

func init() {
	loader.Add("InputsPerconaXtraBackup", func() loader.Plugin { return &InputsPerconaXtraBackup{} })
}
