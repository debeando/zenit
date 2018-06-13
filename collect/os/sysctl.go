package os

import (
  "fmt"
  "gitlab.com/swapbyt3s/zenit/lib"
)

const NR_OPEN string  = "/proc/sys/fs/nr_open"
const FILE_MAX string = "/proc/sys/fs/file-max"

func GetSysLimits(){
  fmt.Printf("os.linux_sysctl_fs.nr_open %d\n", lib.GetValueFromFile(NR_OPEN))
  fmt.Printf("os.linux_sysctl_fs.file_max %d\n", lib.GetValueFromFile(FILE_MAX))
}
