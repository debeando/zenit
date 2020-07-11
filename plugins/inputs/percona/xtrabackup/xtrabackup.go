package process

import (
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
			log.Error("InputsPerconaXtraBackup", map[string]interface{}{"error": err})
		}
	}()

	if !config.File.Inputs.Process.PerconaXtraBackup {
		return
	}

	var a = metrics.Load()
	var pid = common.PGrep("xtrabackup")
	var value = 0

	if pid > 0 {
		value = 1
	}

	log.Debug("InputsPerconaXtraBackup", map[string]interface{}{"xtrabackup": value})

	a.Add(metrics.Metric{
		Key: "process_xtrabackup",
		Tags: []metrics.Tag{
			{"hostname", config.File.General.Hostname},
		},
		Values: []metrics.Value{
			{"xtrabackup", value},
		},
	})
}

func init() {
	inputs.Add("InputsPerconaXtraBackup", func() inputs.Input { return &InputsPerconaXtraBackup{} })
}
