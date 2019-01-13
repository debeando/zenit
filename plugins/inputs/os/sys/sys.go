package sys

import (
	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
	"github.com/swapbyt3s/zenit/plugins/inputs"
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
			{"nr_open", common.GetUInt64FromFile(NR_OPEN)},
			{"file_max", common.GetUInt64FromFile(FILE_MAX)},
		},
	})
}

func init() {
	inputs.Add("InputOSLimits", func() inputs.Input { return &InputOSLimits{} })
}
