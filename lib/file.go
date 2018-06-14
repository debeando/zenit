package lib

import (
  "io/ioutil"
  "os"
  "os/exec"
  "strconv"
  "syscall"
)

func GetValueFromFile(path string) uint64 {
  if _, err := os.Stat(path); err == nil {
    content, err := ioutil.ReadFile(path)
    if err != nil {
      panic(err)
    }

    value, err := strconv.ParseUint(string(content), 10, 64)
    if err != nil {
      panic(err)
    }
    return value
  }

  return 0
}

func PGrep(cmd string) int {
  _, err := exec.Command("/bin/bash", "-c", "/usr/bin/pgrep -x '" + cmd +"' > /dev/null").Output()
  if err != nil {
    if exitError, ok := err.(*exec.ExitError); ok {
      ws := exitError.Sys().(syscall.WaitStatus)
      return ws.ExitStatus()
    }
  }
  return 0
}
