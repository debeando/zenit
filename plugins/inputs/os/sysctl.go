package os

import (
	"github.com/swapbyt3s/zenit/common"
	"github.com/swapbyt3s/zenit/plugins/accumulator"
)

const NR_OPEN string = "/proc/sys/fs/nr_open"
const FILE_MAX string = "/proc/sys/fs/file-max"

func SysLimits() {
	accumulator.Load().Add(accumulator.Metric{
		Key: "os",
		Tags: []accumulator.Tag{
			{"name", "sysctl"},
		},
		Values: []accumulator.Value{{"nr_open", common.GetUInt64FromFile(NR_OPEN)},
			{"file_max", common.GetUInt64FromFile(FILE_MAX)}},
	})
}
