package daemonize

import (
  "fmt"
  "os"
  "os/exec"
  "strconv"
  "syscall"

  "gitlab.com/swapbyt3s/zenit/common"
  "gitlab.com/swapbyt3s/zenit/config"
)

func SavePID(pid int) bool {
  file, err := os.Create(config.General.PIDFile)
  if err != nil {
    return false
  }

  defer file.Close()

  _, err = file.WriteString(strconv.Itoa(pid))
  if err != nil {
    return false
  }

  file.Sync()

  return true
}

func Executable() string {
  ex, err := os.Executable()
  if err != nil {
    panic(err)
  }
  return ex
}

func Build(command string) string {
  return fmt.Sprintf("%s --quiet", command)
}

func Run(command string) int {
  cmd := exec.Command("/bin/bash", "-c", command)
  err := cmd.Start()
  if err != nil {
    panic(err)
  }

  return cmd.Process.Pid
}

func Start() {
  exec := Executable()

  if ! PIDFileExist() {
    cmd := Build(exec)
    pid := Run(cmd)

    if ! SavePID(pid) {
      fmt.Printf("Unable to create PID file: %s\n", config.General.PIDFile)
      os.Exit(1)
    }

    fmt.Printf("Zenit daemon process ID (PID) is %d and is saved in %s\n", pid, config.General.PIDFile)
    os.Exit(0)
  } else {
    fmt.Printf("Zenit already running or %s file exist.\n", config.General.PIDFile)
    os.Exit(1)
  }
}

func Stop() {
  if PIDFileExist() {
    pid := GetPIDFromFile()
    if KillProcess(pid) {
      if RemovePIDFile() {
        os.Exit(0)
      }
    }
  }
  os.Exit(1)
}

func PIDFileExist() bool {
  if _, err := os.Stat(config.General.PIDFile); err != nil {
    return false
  }
  return true
}

func GetPIDFromFile() int {
  return common.GetIntFromFile(config.General.PIDFile)
}

func KillProcess(pid int) bool {
  if err := syscall.Kill(pid, syscall.SIGKILL); err != nil {
    return false
  }
  return true
}

func RemovePIDFile() bool {
  if err := os.Remove(config.General.PIDFile); err != nil {
    return false
  }
  return true
}
