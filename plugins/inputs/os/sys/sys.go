package sys

import (
	"fmt"

	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/common/log"
	"github.com/swapbyt3s/zenit/config"
	"github.com/swapbyt3s/zenit/plugins/inputs"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

const NR_OPEN string = "/proc/sys/fs/nr_open"
const FILE_MAX string = "/proc/sys/fs/file-max"

type InputOSLimits struct{}

func (l *InputOSLimits) Collect() {
	defer func() {
		if err := recover(); err != nil {
			log.Debug(fmt.Sprintf("Plugin - InputOSLimits - Panic (code %d) has been recover from somewhere.\n", err))
		}
	}()

	if !config.File.OS.Inputs.Limits {
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
