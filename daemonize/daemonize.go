package daemonize

import (
  "fmt"
  "log"
  "os"
  "os/exec"
  "strconv"
  "strings"
)

var PIDFile = "/tmp/zenit.pid"

func savePID(pid int) {
  file, err := os.Create(PIDFile)
  if err != nil {
    log.Printf("Unable to create pid file : %v\n", err)
    os.Exit(1)
  }

  defer file.Close()

  _, err = file.WriteString(strconv.Itoa(pid))
  if err != nil {
    log.Printf("Unable to create pid file : %v\n", err)
    os.Exit(1)
  }

  file.Sync()
}

func getExecutable() string {
  ex, err := os.Executable()
  if err != nil {
    panic(err)
  }
  return ex
}

func getArgs() string {
  cmd := strings.Join(os.Args[1:], " ")
  cmd  = strings.Replace(cmd, "--daemonize", "", -1)
  cmd  = strings.Replace(cmd, "-daemonize", "", -1)
  cmd  = strings.TrimSpace(cmd)

  return cmd
}

func getCommand(command string, args string) string {
  return command + " " + args
}

func runCommand(command string) int {
  cmd := exec.Command("/bin/bash", "-c", command)
  err := cmd.Start()
  if err != nil {
    panic(err)
  }

  return cmd.Process.Pid
}

func Start() {
  // Check if daemon already running.
  if _, err := os.Stat(PIDFile); err == nil {
    fmt.Printf("Already running or %s file exist.\n", PIDFile)
    os.Exit(1)
  }

  exec := getExecutable()
  args := getArgs()
  cmd  := getCommand(exec, args)
  pid  := runCommand(cmd)

  fmt.Println("Daemon process ID is: ", pid)

  savePID(pid)

  os.Exit(0)
}
