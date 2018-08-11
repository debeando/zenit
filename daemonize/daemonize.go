package daemonize

import (
  "fmt"
  "os"
  "os/exec"
  "syscall"

  "github.com/swapbyt3s/zenit/common"
  "github.com/swapbyt3s/zenit/common/file"
  "github.com/swapbyt3s/zenit/config"
)

func Run(command string) int {
  cmd := exec.Command("/bin/bash", "-c", command)
  err := cmd.Start()
  if err != nil {
    panic(err)
  }

  return cmd.Process.Pid
}

func Start() {
  if ! file.Exist(config.General.PIDFile) {
    exec, _ := os.Executable()
    cmd  := fmt.Sprintf("%s --quiet", exec)
    pid  := Run(cmd)

    if file.Create(config.General.PIDFile) {
      if file.Write(config.General.PIDFile, common.IntToString(pid)) {
        fmt.Printf("Zenit daemon process ID (PID) is %d and is saved in %s\n", pid, config.General.PIDFile)
        os.Exit(0)
      }
    }

    fmt.Printf("Unable to create PID file: %s\n", config.General.PIDFile)
    os.Exit(1)
  } else {
    fmt.Printf("Zenit already running or %s file exist.\n", config.General.PIDFile)
    os.Exit(1)
  }
}

func Stop() {
  if file.Exist(config.General.PIDFile) {
    pid := common.GetIntFromFile(config.General.PIDFile)
    if Kill(pid) {
      if file.Delete(config.General.PIDFile) {
        os.Exit(0)
      }
    }
  }
  os.Exit(1)
}

func Kill(pid int) bool {
  pgid, err := syscall.Getpgid(pid)
  if err == nil {
    if err := syscall.Kill(-pgid, syscall.SIGKILL); err != nil {
      return false
    }
  }
  return true
}
