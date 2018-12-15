package sys

import (
	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/loader"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

const NR_OPEN string = "/proc/sys/fs/nr_open"
const FILE_MAX string = "/proc/sys/fs/file-max"

type InputOSLimits struct {}

func (l *InputOSLimits) Collect() {
	if ! config.File.OS.Inputs.Limits {
		return
	}

	metrics.Load().Add(metrics.Metric{
		Key: "zenit_os",
		Tags: []metrics.Tag{
			{"name", "sysctl"},
		},
		Values: []metrics.Value{
			{"nr_open", uint(common.GetUInt64FromFile(NR_OPEN))},
			{"file_max", uint(common.GetUInt64FromFile(FILE_MAX))},
		},
	})
}

func init() {
	loader.Add("InputOSLimits", func() loader.Plugin { return &InputOSLimits{} })
}
