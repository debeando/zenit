// TODO: Read from zenit.yaml the list of process to check.
package process

import (
	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/loader"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

type InputsPerconaXtraBackup struct {}

func (l *InputsPerconaXtraBackup) Collect() {
	if ! config.File.Process.Inputs.PerconaXtraBackup {
		return
	}

	metrics.Load().Add(metrics.Metric{
		Key: "zenit_os",
		Tags: []metrics.Tag{{"system", "linux"},
			{"process", "xtrabackup"}},
		Values: common.PGrep("xtrabackup") ^ 1,
	})
}

func init() {
	loader.Add("InputsPerconaXtraBackup", func() loader.Plugin { return &InputsPerconaXtraBackup{} })
}
