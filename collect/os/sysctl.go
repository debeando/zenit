package os

import (
  "fmt"
  "io/ioutil"
  "os"
  "strconv"
  "strings"
)

const NR_OPEN string  = "/proc/sys/fs/nr_open"
const FILE_MAX string = "/proc/sys/fs/file-max"

func GetSysLimits(){
  fmt.Printf("os.linux_sysctl_fs.nr_open %d\n", GetValueFromFile(NR_OPEN))
  fmt.Printf("os.linux_sysctl_fs.file_max %d\n", GetValueFromFile(FILE_MAX))
}

func GetValueFromFile(path string) uint64 {
  if _, err := os.Stat(path); err == nil {
    content, err := ioutil.ReadFile(path)
    if err != nil {
      panic(err)
    }

    value, err := strconv.ParseUint(strings.Trim(string(content), "\n"), 10, 64)
    if err != nil {
      panic(err)
    }
    return value
  }

  return 0
}