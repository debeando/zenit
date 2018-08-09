package common

import (
  "bufio"
  "fmt"
  "log"
  "os"
  "os/exec"
)

func Tail(path string, channel chan<- string) {
  if _, err := os.Stat(path); err != nil {
    log.Printf("E! - Tail - File not exist: %s\n", path)
    os.Exit(1)
  }

  cmd := exec.Command("/usr/bin/tail", "-n", "0", "-f", path)

  cmdReader, err := cmd.StdoutPipe()
  if err != nil {
    fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
    os.Exit(1)
  }

  scanner := bufio.NewScanner(cmdReader)

  go func() {
    defer close(channel)

    for scanner.Scan() {
      channel <- scanner.Text()
    }
  }()

  err = cmd.Start()
  if err != nil {
    fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
    os.Exit(1)
  }

  err = cmd.Wait()
  if err != nil {
    fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
    os.Exit(1)
  }
}
