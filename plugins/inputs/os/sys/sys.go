package sys

import (
	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/plugins/lists/metrics"
)

const NR_OPEN string = "/proc/sys/fs/nr_open"
const FILE_MAX string = "/proc/sys/fs/file-max"

func Collect() {
	metrics.Load().Add(metrics.Metric{
		Key: "os",
		Tags: []metrics.Tag{
			{"name", "sysctl"},
		},
		Values: []metrics.Value{{"nr_open", common.GetUInt64FromFile(NR_OPEN)},
			{"file_max", common.GetUInt64FromFile(FILE_MAX)}},
	})
}
