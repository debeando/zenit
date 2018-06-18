package os

import (
  "fmt"
  "gitlab.com/swapbyt3s/zenit/common"
)

const NR_OPEN string  = "/proc/sys/fs/nr_open"
const FILE_MAX string = "/proc/sys/fs/file-max"

func GatherSysLimits(){
  fmt.Printf("os.linux_sysctl_fs.nr_open %d\n", common.ValueFromFile(NR_OPEN))
  fmt.Printf("os.linux_sysctl_fs.file_max %d\n", common.ValueFromFile(FILE_MAX))
}
