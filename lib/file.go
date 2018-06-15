package lib

import (
  "os/exec"
  "syscall"
)

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
