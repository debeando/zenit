package lib

import (
  "io/ioutil"
  "os"
  "os/exec"
  "strconv"
  "strings"
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

func ReadFile(path string) (lines []string) {
  if _, err := os.Stat(path); err == nil {
    contents, err := ioutil.ReadFile(path)
    if err != nil {
      panic(err)
    }

    lines = strings.Split(string(contents), "\n")
  }
  return
}

func ValueFromFile(path string) uint64 {
  lines := ReadFile(path)
  if len(lines) > 0 {
    return StringToUInt64(lines[0])
  }
  return 0
}

func StringToUInt64(value string) uint64 {
  i, err := strconv.ParseUint(strings.TrimSpace(value), 10, 64)
  if err != nil {
    panic(err)
  }
  return i
}
