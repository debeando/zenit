package daemonize

import (
  "fmt"
  "os"
  "os/exec"
  "strconv"
  "strings"
  "gitlab.com/swapbyt3s/zenit/common"
)

const PIDFile = "/tmp/zenit-%s.pid"

func SavePID(filename string, pid int) bool {
  file, err := os.Create(filename)
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

func Args(args []string) string {
  a := strings.Join(args[1:], " ")
  a  = strings.Replace(a, "--daemonize", "", -1)
  a  = strings.Replace(a, "-daemonize", "", -1)
  a  = strings.TrimSpace(a)

  return a
}

func Build(command string, args string) string {
  return fmt.Sprintf("%s %s", command, args)
}

func Run(command string) int {
  cmd := exec.Command("/bin/bash", "-c", command)
  err := cmd.Start()
  if err != nil {
    panic(err)
  }

  return cmd.Process.Pid
}

func GetPIDFileName(args string) string {
  return fmt.Sprintf(PIDFile, common.MD5(args))
}

func Start() {
  exec := Executable()
  args := Args(os.Args)
  file := GetPIDFileName(args)

  if ! PIDFileExist(file) {
    cmd := Build(exec, args)
    pid := Run(cmd)

    if ! SavePID(file, pid) {
      fmt.Printf("Unable to create PID file: %s\n", file)
      os.Exit(1)
    }

    fmt.Printf("Zenit daemon process ID (PID) is %d and is saved in %s\n", pid, file)
    os.Exit(0)
  } else {
    fmt.Printf("Zenit already running or %s file exist.\n", file)
    os.Exit(1)
  }
}

func Stop() {
  exec := Executable()
  args := Args(os.Args)
  file := GetPIDFileName(args)

  if PIDFileExist(file) {

  }
}

func PIDFileExist(file string) bool {
  if _, err := os.Stat(file); err == nil {
    return true
  }
  return false
}

func GetPIDFromFile() {

}

func KillProcess() {

}

func RemovePIDFile{

}
