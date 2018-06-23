package os

import (
  "gitlab.com/swapbyt3s/zenit/common"
  "gitlab.com/swapbyt3s/zenit/output"
)

const NR_OPEN string  = "/proc/sys/fs/nr_open"
const FILE_MAX string = "/proc/sys/fs/file-max"

func GatherSysLimits(){
  output.Load().AddItem(output.Metric{
    Key: "os",
    Tags: []output.Tag{output.Tag{"system", "linux"},
                       output.Tag{"setting", "sysctl"}},
    Values: []output.Value{output.Value{"nr_open", common.ValueFromFile(NR_OPEN)},
                           output.Value{"file_max", common.ValueFromFile(FILE_MAX)}},
  })
}
