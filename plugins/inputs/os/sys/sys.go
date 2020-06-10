package sys

import (
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
			log.Error("InputOSLimits", map[string]interface{}{"error": err})
		}
	}()

	if !config.File.Inputs.OS.Limits {
		return
	}

	log.Debug("InputOSLimits", map[string]interface{}{
		"nr_open": common.GetInt64FromFile(NR_OPEN),
		"file_max": common.GetInt64FromFile(FILE_MAX),
	})

	var a = metrics.Load()

	a.Add(metrics.Metric{
		Key: "os_sys",
		Tags: []metrics.Tag{
			{"hostname", config.File.General.Hostname},
		},
		Values: []metrics.Value{
			{"nr_open", common.GetInt64FromFile(NR_OPEN)},
			{"file_max", common.GetInt64FromFile(FILE_MAX)},
		},
	})
}

func init() {
	inputs.Add("InputOSLimits", func() inputs.Input { return &InputOSLimits{} })
}
