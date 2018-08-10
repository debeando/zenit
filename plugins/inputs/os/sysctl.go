package os

import (
  "github.com/swapbyt3s/zenit/common"
  "github.com/swapbyt3s/zenit/plugins/accumulator"
)

const NR_OPEN string  = "/proc/sys/fs/nr_open"
const FILE_MAX string = "/proc/sys/fs/file-max"

func SysLimits(){
  accumulator.Load().AddItem(accumulator.Metric{
    Key: "os",
    Tags: []accumulator.Tag{accumulator.Tag{"system", "linux"},
                            accumulator.Tag{"setting", "sysctl"}},
    Values: []accumulator.Value{accumulator.Value{"nr_open", common.GetUInt64FromFile(NR_OPEN)},
                                accumulator.Value{"file_max", common.GetUInt64FromFile(FILE_MAX)}},
  })
}
