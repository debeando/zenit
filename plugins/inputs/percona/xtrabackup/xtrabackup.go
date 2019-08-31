package process

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/inputs"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

type InputsPerconaXtraBackup struct{}

func (l *InputsPerconaXtraBackup) Collect() {
	defer func() {
		if err := recover(); err != nil {
			log.Debug(fmt.Sprintf("Plugin - InputsPerconaXtraBackup - Panic (code %d) has been recover from somewhere.\n", err))
		}
	}()

	if !config.File.Inputs.Process.PerconaXtraBackup {
		return
	}

	var pid = common.PGrep("xtrabackup")
	var value = 0

	if pid > 0 {
		value = 1
	}

	metrics.Load().Add(metrics.Metric{
		Key: "process",
		Tags: []metrics.Tag{},
		Values: []metrics.Value{
			{ "xtrabackup", value},
		},
	})

	log.Debug(fmt.Sprintf("Plugin - InputsPerconaXtraBackup - %d", value))
}

func init() {
	inputs.Add("InputsPerconaXtraBackup", func() inputs.Input { return &InputsPerconaXtraBackup{} })
}
